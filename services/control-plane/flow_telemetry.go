package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"net"
	"sync"
)

const (
	maxFlowPacketSize = 65535
	maxTemplatesPerExporter = 4096
	maxRecordsPerPacket = 4096
)

var (
	errFlowPacketTooLarge = errors.New("flow packet exceeds maximum size")
	errFlowPacketTruncated = errors.New("flow packet is truncated")
	errFlowUnsupportedVersion = errors.New("unsupported flow version")
	errFlowTemplateMissing = errors.New("flow template missing")
	errFlowInvalidCount = errors.New("invalid flow record count")
)

type FlowExporterKey struct { Address string; Protocol string }
type FlowTemplateField struct { IE uint16; Length uint16; Enterprise uint32 }
type FlowTemplate struct { ID uint16; Fields []FlowTemplateField }
type FlowTemplateCache struct { mu sync.RWMutex; templates map[FlowExporterKey]map[uint16]FlowTemplate }

func NewFlowTemplateCache() *FlowTemplateCache { return &FlowTemplateCache{templates: make(map[FlowExporterKey]map[uint16]FlowTemplate)} }
func (c *FlowTemplateCache) Put(key FlowExporterKey, t FlowTemplate) error { if t.ID<256||len(t.Fields)==0||len(t.Fields)>64{return errors.New("invalid flow template")}; c.mu.Lock();defer c.mu.Unlock();m:=c.templates[key];if m==nil{m=make(map[uint16]FlowTemplate);c.templates[key]=m};if _,exists:=m[t.ID];!exists&&len(m)>=maxTemplatesPerExporter{return errors.New("template cache limit reached")};m[t.ID]=t;return nil }
func (c *FlowTemplateCache) Get(key FlowExporterKey,id uint16)(FlowTemplate,bool){c.mu.RLock();defer c.mu.RUnlock();m:=c.templates[key];if m==nil{return FlowTemplate{},false};t,ok:=m[id];return t,ok}
func readU16(b []byte,off *int)(uint16,error){if *off+2>len(b){return 0,errFlowPacketTruncated};v:=binary.BigEndian.Uint16(b[*off:*off+2]);*off+=2;return v,nil}

func DecodeNetFlowV5(packet []byte, exporter string)([]FlowRecord,error){if len(packet)>maxFlowPacketSize{return nil,errFlowPacketTooLarge};if len(packet)<24{return nil,errFlowPacketTruncated};if binary.BigEndian.Uint16(packet[:2])!=5{return nil,errFlowUnsupportedVersion};count:=int(binary.BigEndian.Uint16(packet[2:4]));if count>30||count>maxRecordsPerPacket{return nil,errFlowInvalidCount};if len(packet)<24+count*48{return nil,errFlowPacketTruncated};out:=make([]FlowRecord,0,count);for i:=0;i<count;i++{base:=24+i*48;src:=net.IP(packet[base:base+4]).String();dst:=net.IP(packet[base+4:base+8]).String();sport:=binary.BigEndian.Uint16(packet[base+32:base+34]);dport:=binary.BigEndian.Uint16(packet[base+34:base+36]);proto:=packet[base+38];packets:=uint64(binary.BigEndian.Uint32(packet[base+16:base+20]));bytes:=uint64(binary.BigEndian.Uint32(packet[base+20:base+24]));out=append(out,FlowRecord{ExporterID:exporter,Version:5,SourceIP:src,DestinationIP:dst,SourcePort:sport,DestinationPort:dport,Protocol:proto,Bytes:bytes,Packets:packets,SamplingRate:1})};return out,nil}

func ParseNetFlowV9Templates(packet []byte,key FlowExporterKey,cache *FlowTemplateCache)error{if len(packet)>maxFlowPacketSize{return errFlowPacketTooLarge};if len(packet)<20{return errFlowPacketTruncated};if binary.BigEndian.Uint16(packet[:2])!=9{return errFlowUnsupportedVersion};count:=int(binary.BigEndian.Uint16(packet[2:4]));if count>maxRecordsPerPacket{return errFlowInvalidCount};off:=20;for off+4<=len(packet)&&count>0{setID:=binary.BigEndian.Uint16(packet[off:off+2]);setLen:=int(binary.BigEndian.Uint16(packet[off+2:off+4]));if setLen<4||off+setLen>len(packet){return errFlowPacketTruncated};if setID==0{p:=off+4;for p+4<=off+setLen{id,err:=readU16(packet,&p);if err!=nil{return err};fields,err:=readU16(packet,&p);if err!=nil{return err};if id<256||fields==0||fields>64{return errors.New("invalid v9 template")};f:=make([]FlowTemplateField,0,fields);for j:=0;j<int(fields);j++{ie,err:=readU16(packet,&p);if err!=nil{return err};ln,err:=readU16(packet,&p);if err!=nil{return err};if ln==0{return errors.New("invalid v9 template field length")};f=append(f,FlowTemplateField{IE:ie,Length:ln})};if err:=cache.Put(key,FlowTemplate{ID:id,Fields:f});err!=nil{return err}}};off+=setLen;count--};return nil}
func decodeUint(b []byte)uint64{var v uint64;for _,x:=range b{v=(v<<8)|uint64(x)};return v}

func DecodeIPFIXDataSet(packet []byte,key FlowExporterKey,templateID uint16,cache *FlowTemplateCache,exporter string)([]FlowRecord,error){if len(packet)>maxFlowPacketSize{return nil,errFlowPacketTooLarge};t,ok:=cache.Get(key,templateID);if !ok{return nil,errFlowTemplateMissing};recordLen:=0;for _,f:=range t.Fields{if f.Length==65535{return nil,errors.New("unsupported variable-length flow field")};recordLen+=int(f.Length)};if recordLen<=0{return nil,errors.New("invalid template record length")};if len(packet)%recordLen!=0{return nil,errFlowPacketTruncated};count:=len(packet)/recordLen;if count>maxRecordsPerPacket{return nil,errFlowInvalidCount};out:=make([]FlowRecord,0,count);for n:=0;n<count;n++{var r FlowRecord;r.ExporterID=exporter;r.Version=10;off:=n*recordLen;for _,f:=range t.Fields{v:=packet[off:off+int(f.Length)];off+=int(f.Length);if f.Enterprise!=0{continue};switch f.IE{case 1:r.Bytes=decodeUint(v);case 2:r.Packets=decodeUint(v);case 4:if len(v)>0{r.Protocol=v[len(v)-1]};case 7:if len(v)<=2{r.SourcePort=uint16(decodeUint(v))};case 8:if len(v)==4{r.SourceIP=net.IP(v).String()};case 11:if len(v)<=2{r.DestinationPort=uint16(decodeUint(v))};case 12:if len(v)==4{r.DestinationIP=net.IP(v).String()};case 27:if len(v)==16{r.SourceIP=net.IP(v).String()};case 28:if len(v)==16{r.DestinationIP=net.IP(v).String()};case 34:r.SamplingRate=uint32(decodeUint(v))}};if r.SamplingRate==0{r.SamplingRate=1};out=append(out,r)};return out,nil}
func FlowRecordFingerprint(r FlowRecord)uint64{h:=fnv.New64a();fmt.Fprintf(h,"%s|%s|%d|%d|%d|%d|%d|%d",r.SourceIP,r.DestinationIP,r.SourcePort,r.DestinationPort,r.Protocol,r.Bytes,r.Packets,r.SamplingRate);return h.Sum64()}
func NormalizeSampledCounters(r FlowRecord)FlowRecord{if r.SamplingRate<=1{if r.SamplingRate==0{r.SamplingRate=1};return r};rate:=uint64(r.SamplingRate);max:=^uint64(0);if r.Bytes>max/rate{r.Bytes=max}else{r.Bytes*=rate};if r.Packets>max/rate{r.Packets=max}else{r.Packets*=rate};return r}

type FlowTelemetryCollector struct{Templates *FlowTemplateCache}
func NewFlowTelemetryCollector()*FlowTelemetryCollector{return &FlowTelemetryCollector{Templates:NewFlowTemplateCache()}}
func(c *FlowTelemetryCollector)Protocol()string{return "netflow-ipfix"}
func(c *FlowTelemetryCollector)Collect(ctx context.Context,packet []byte)([]FlowRecord,error){_ = ctx;if len(packet)<2{return nil,errFlowPacketTruncated};switch binary.BigEndian.Uint16(packet[:2]){case 5:return DecodeNetFlowV5(packet,"unknown");case 9:return nil,errors.New("netflow v9 requires exporter identity and template-aware decoding");case 10:return nil,errors.New("ipfix requires exporter identity and template-aware decoding");default:return nil,errFlowUnsupportedVersion}}
