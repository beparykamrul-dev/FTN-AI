package observability

import("context";"strings";"time")
type Collector interface{Name() string;Collect(context.Context)([]TrafficSample,error)}
type InterfaceCounter struct{Interface string `json:"interface"`;RxBytes uint64 `json:"rx_bytes"`;TxBytes uint64 `json:"tx_bytes"`;RxPackets uint64 `json:"rx_packets"`;TxPackets uint64 `json:"tx_packets"`;At time.Time `json:"at"`}
type CounterSource interface{ReadCounters(context.Context)([]InterfaceCounter,error)}
type InterfaceCollector struct{Source CounterSource}
func(c InterfaceCollector)Name()string{return "interface-counters"}
func(c InterfaceCollector)Collect(ctx context.Context)([]TrafficSample,error){if ctx==nil{return nil,context.Canceled};if c.Source==nil{return nil,context.Canceled};if err:=ctx.Err();err!=nil{return nil,err};counters,err:=c.Source.ReadCounters(ctx);if err!=nil{return nil,err};out:=make([]TrafficSample,0,len(counters));for _,v:=range counters{iface:=strings.TrimSpace(v.Interface);if iface==""{continue};at:=v.At;if at.IsZero(){at=time.Now().UTC()};out=append(out,TrafficSample{Interface:iface,Bytes:v.RxBytes+v.TxBytes,Packets:v.RxPackets+v.TxPackets,At:at})};return out,nil}
