package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

type NFCAPDFlowAdapter struct{}

type nfcapdColumns struct { src, dst, sport, dport, proto, bytes, packets, sampling, first, last int }

func (NFCAPDFlowAdapter) Name() string { return "nfcapd-csv" }
func normalizeNFCAPDColumn(s string) string { s=strings.ToLower(strings.TrimSpace(s)); s=strings.Trim(s,"\"'"); return strings.ReplaceAll(s,"_","") }
func findNFCAPDColumn(header []string, aliases ...string) int { want:=make(map[string]struct{},len(aliases));for _,a:=range aliases{want[normalizeNFCAPDColumn(a)]=struct{}{}};for i,n:=range header{if _,ok:=want[normalizeNFCAPDColumn(n)];ok{return i}};return -1 }
func parseNFCAPDUint(rec []string, idx int, field string, bits int) (uint64,error) { if idx<0||idx>=len(rec){return 0,fmt.Errorf("%s_required",field)};v:=strings.TrimSpace(rec[idx]);if v==""{return 0,fmt.Errorf("%s_required",field)};n,e:=strconv.ParseUint(v,10,bits);if e!=nil{return 0,fmt.Errorf("%s: %w",field,e)};return n,nil }
func sumNFCAPD(a,b uint64)(uint64,error){if math.MaxUint64-a<b{return 0,errors.New("counter_overflow")};return a+b,nil}

func parseNFCAPDTimestamp(value string) (time.Time,error) {
	value=strings.TrimSpace(strings.Trim(value,"\"'"));if value==""{return time.Time{},errors.New("timestamp_required")}
	for _,layout:=range []string{time.RFC3339Nano,time.RFC3339,"2006-01-02 15:04:05.999999999 -0700 MST","2006-01-02 15:04:05"}{
		if t,e:=time.Parse(layout,value);e==nil{return t.UTC(),nil}
	}
	if n,e:=strconv.ParseInt(value,10,64);e==nil { if n>1_000_000_000_000 { return time.UnixMilli(n).UTC(),nil }; return time.Unix(n,0).UTC(),nil }
	return time.Time{},fmt.Errorf("invalid_timestamp: %s",value)
}

func (a NFCAPDFlowAdapter) decode(ctx context.Context,r io.Reader,exporter string,version uint16)([]FlowObservation,error){
	if r==nil{return nil,errors.New("flow_reader_required")};if ctx==nil{ctx=context.Background()};exporter=strings.TrimSpace(exporter);if exporter==""||net.ParseIP(exporter)==nil{return nil,errors.New("invalid_exporter")};if version!=5&&version!=9&&version!=10{return nil,errors.New("unsupported_flow_version")}
	cr:=csv.NewReader(bufio.NewReader(r));cr.TrimLeadingSpace=true;cr.FieldsPerRecord=-1;out:=make([]FlowObservation,0,128);var columns *nfcapdColumns
	for {select{case <-ctx.Done():return out,ctx.Err();default:{}};rec,e:=cr.Read();if e==io.EOF{break};if e!=nil{return nil,fmt.Errorf("read_nfcapd_csv: %w",e)};if len(rec)==0{continue};for i:=range rec{rec[i]=strings.TrimSpace(rec[i])}
		if columns==nil&&len(rec)>=8&&(strings.EqualFold(rec[0],"ts")||strings.EqualFold(rec[0],"first")||strings.EqualFold(rec[0],"srcip")||strings.EqualFold(rec[0],"sa")){c:=&nfcapdColumns{src:findNFCAPDColumn(rec,"srcip","sa","sourceip","sourceipv4address","sourceipv6address"),dst:findNFCAPDColumn(rec,"dstip","da","destinationip","destinationipv4address","destinationipv6address"),sport:findNFCAPDColumn(rec,"srcport","sp","spt","sourcetransportport"),dport:findNFCAPDColumn(rec,"dstport","dp","dpt","destinationtransportport"),proto:findNFCAPDColumn(rec,"proto","pr","protocol","protocolidentifier"),bytes:findNFCAPDColumn(rec,"bytes","octettotalcount","ibytes","obytes"),packets:findNFCAPDColumn(rec,"packets","packettotalcount","ipackets","opackets"),sampling:findNFCAPDColumn(rec,"sampling","samplinginterval"),first:findNFCAPDColumn(rec,"first","ts","firstseen","firstseenat"),last:findNFCAPDColumn(rec,"last","te","lastseen","lastseenat")};if c.src>=0&&c.dst>=0&&c.sport>=0&&c.dport>=0&&c.proto>=0&&c.bytes>=0&&c.packets>=0{columns=c;continue}}
		if len(rec)>=1&&(strings.EqualFold(rec[0],"ts")||strings.EqualFold(rec[0],"first")||strings.EqualFold(rec[0],"srcip")||strings.EqualFold(rec[0],"sa")){continue}
		si,di,spi,dpi,pi,bi,ki,ri:=0,1,2,3,4,5,6,7;fi,li:=-1,-1;if columns!=nil{si,di,spi,dpi,pi,bi,ki,ri=columns.src,columns.dst,columns.sport,columns.dport,columns.proto,columns.bytes,columns.packets,columns.sampling;fi,li=columns.first,columns.last};if len(rec)<=bi||len(rec)<=ki{return nil,errors.New("nfcapd_record_too_short")};if si<0||di<0||si>=len(rec)||di>=len(rec){return nil,errors.New("flow_source_destination_required")};src,dst:=rec[si],rec[di];if net.ParseIP(src)==nil||net.ParseIP(dst)==nil{return nil,errors.New("flow_source_destination_required")}
		sport,e:=parseNFCAPDUint(rec,spi,"source_port",16);if e!=nil{return nil,e};dport,e:=parseNFCAPDUint(rec,dpi,"destination_port",16);if e!=nil{return nil,e};proto,e:=parseNFCAPDUint(rec,pi,"protocol",8);if e!=nil{return nil,e};bytes,e:=parseNFCAPDUint(rec,bi,"bytes",64);if e!=nil{return nil,e};packets,e:=parseNFCAPDUint(rec,ki,"packets",64);if e!=nil{return nil,e};rate:=uint64(1);if ri>=0&&ri<len(rec)&&rec[ri]!=""{rate,e=parseNFCAPDUint(rec,ri,"sampling_rate",32);if e!=nil{return nil,e};if rate==0{rate=1}}
		first,last:=time.Time{},time.Time{};if fi>=0&&fi<len(rec)&&rec[fi]!=""{first,e=parseNFCAPDTimestamp(rec[fi]);if e!=nil{return nil,e}};if li>=0&&li<len(rec)&&rec[li]!=""{last,e=parseNFCAPDTimestamp(rec[li]);if e!=nil{return nil,e}};observed:=last;if observed.IsZero(){observed=first};if observed.IsZero(){return nil,errors.New("observed_timestamp_required")}
		o,e:=normalizeFlowObservation(FlowRecord{ExporterID:exporter,Version:version,SourceIP:src,DestinationIP:dst,SourcePort:uint16(sport),DestinationPort:uint16(dport),Protocol:uint8(proto),Bytes:bytes,Packets:packets,SamplingRate:uint32(rate)},observed,first,last);if e!=nil{return nil,e};out=append(out,o);if len(out)>maxSiLKBatchSize{return nil,errors.New("flow_batch_limit")}
	};return out,nil
}

// DecodeObservations preserves exporter timestamps for durable ingestion.
func (a NFCAPDFlowAdapter) DecodeObservations(ctx context.Context,r io.Reader,exporter string,version uint16)([]FlowObservation,error){return a.decode(ctx,r,exporter,version)}

func (a NFCAPDFlowAdapter) Decode(ctx context.Context,r io.Reader,exporter string,version uint16)([]FlowRecord,error){obs,e:=a.decode(ctx,r,exporter,version);if e!=nil{return nil,e};out:=make([]FlowRecord,0,len(obs));for _,o:=range obs{out=append(out,o.FlowRecord)};return out,nil}
