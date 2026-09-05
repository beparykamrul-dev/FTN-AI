package observability

import (
	"sort"
	"strings"
)

type TrafficAggregate struct { Interface string `json:"interface"`; Protocol string `json:"protocol,omitempty"`; App string `json:"application,omitempty"`; Bytes uint64 `json:"bytes"`; Packets uint64 `json:"packets"` }
func Aggregate(samples []TrafficSample) []TrafficAggregate { m:=make(map[string]*TrafficAggregate); for _,s:=range samples{s.Interface=strings.TrimSpace(s.Interface);s.Protocol=strings.ToLower(strings.TrimSpace(s.Protocol));s.App=strings.TrimSpace(s.App);if s.Interface==""{continue};key:=s.Interface+"|"+s.Protocol+"|"+s.App;v:=m[key];if v==nil{v=&TrafficAggregate{Interface:s.Interface,Protocol:s.Protocol,App:s.App};m[key]=v};v.Bytes+=s.Bytes;v.Packets+=s.Packets};out:=make([]TrafficAggregate,0,len(m));for _,v:=range m{out=append(out,*v)};sort.SliceStable(out,func(i,j int)bool{if out[i].Bytes!=out[j].Bytes{return out[i].Bytes>out[j].Bytes};if out[i].Packets!=out[j].Packets{return out[i].Packets>out[j].Packets};if out[i].Interface!=out[j].Interface{return out[i].Interface<out[j].Interface};if out[i].Protocol!=out[j].Protocol{return out[i].Protocol<out[j].Protocol};return out[i].App<out[j].App});return out}
