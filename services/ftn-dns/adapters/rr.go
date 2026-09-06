package adapters

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type RR struct { Name string; Type uint16; Class uint16; TTL uint32; Value string }

func ParseAnswers(msg []byte) ([]RR, error) {
	if len(msg)<12{return nil,fmt.Errorf("dns message too short")};qd,an:=int(binary.BigEndian.Uint16(msg[4:6])),int(binary.BigEndian.Uint16(msg[6:8]));if qd>1024||an>65535{return nil,fmt.Errorf("dns section count too large")};off:=12
	for i:=0;i<qd;i++{var err error;off,err=skipName(msg,off);if err!=nil||off+4>len(msg){return nil,fmt.Errorf("invalid question section")};off+=4}
	answers:=make([]RR,0,an);for i:=0;i<an;i++{name,next,err:=readName(msg,off);if err!=nil||next+10>len(msg){return nil,fmt.Errorf("invalid answer name")};typ:=binary.BigEndian.Uint16(msg[next:next+2]);class:=binary.BigEndian.Uint16(msg[next+2:next+4]);ttl:=binary.BigEndian.Uint32(msg[next+4:next+8]);rdlen:=int(binary.BigEndian.Uint16(msg[next+8:next+10]));rstart:=next+10;if rstart+rdlen>len(msg){return nil,fmt.Errorf("invalid rdata length")};value:=decodeRData(msg,typ,rstart,rdlen);answers=append(answers,RR{Name:name,Type:typ,Class:class,TTL:ttl,Value:value});off=rstart+rdlen}
	return answers,nil
}
func skipName(msg []byte,off int)(int,error){_,next,err:=readName(msg,off);return next,err}
func readName(msg []byte,off int)(string,int,error){if off<0||off>=len(msg){return "",0,fmt.Errorf("name offset out of range")};labels:=[]string{};cur,next:=off,off;jumped:=false;seen:=map[int]bool{};for steps:=0;steps<256;steps++{if cur>=len(msg){return "",0,fmt.Errorf("name truncated")};if seen[cur]{return "",0,fmt.Errorf("name compression loop")};seen[cur]=true;ln:=int(msg[cur]);if ln&0xc0==0xc0{if cur+1>=len(msg){return "",0,fmt.Errorf("pointer truncated")};ptr:=((ln&0x3f)<<8)|int(msg[cur+1]);if ptr>=len(msg){return "",0,fmt.Errorf("pointer out of range")};if !jumped{next=cur+2};cur=ptr;jumped=true;continue};if ln==0{if !jumped{next=cur+1};break};if ln>63||cur+1+ln>len(msg){return "",0,fmt.Errorf("invalid label")};label:=string(msg[cur+1:cur+1+ln]);if len(label)>63{return "",0,fmt.Errorf("invalid label")};labels=append(labels,label);cur+=1+ln;if len(labels)>127{return "",0,fmt.Errorf("too many labels")}};if cur<len(msg)&&len(labels)>=127{return "",0,fmt.Errorf("name too long")};return strings.Join(labels,"."),next,nil}
func decodeRData(msg []byte,typ uint16,start,length int)string{switch typ{case 1:if length==net.IPv4len{return net.IP(msg[start:start+length]).String()};case 28:if length==net.IPv6len{return net.IP(msg[start:start+length]).String()};case 16:if length>1&&int(msg[start])<=length-1{return string(msg[start+1:start+1+int(msg[start])])};case 2,5,12:if name,_,err:=readName(msg,start);err==nil{return name}};return fmt.Sprintf("rdata:%d-bytes",length)}
