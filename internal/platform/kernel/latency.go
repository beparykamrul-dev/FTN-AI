package kernel

import (
	"sort"
	"strings"
	"sync"
	"time"
)
type EndpointSample struct{Endpoint string;RTT time.Duration;Loss float64;Healthy bool;ObservedAt time.Time}
type LatencyRouter struct{mu sync.RWMutex;samples map[string]EndpointSample}
func NewLatencyRouter()*LatencyRouter{return &LatencyRouter{samples:make(map[string]EndpointSample)}}
func(r *LatencyRouter)Observe(s EndpointSample){if r==nil{return};s.Endpoint=strings.TrimSpace(s.Endpoint);if s.Endpoint==""||s.RTT<0||s.Loss<0||s.Loss>1{return};if s.ObservedAt.IsZero(){s.ObservedAt=time.Now().UTC()};r.mu.Lock();r.samples[s.Endpoint]=s;r.mu.Unlock()}
func(r *LatencyRouter)Best()(EndpointSample,bool){if r==nil{return EndpointSample{},false};r.mu.RLock();items:=make([]EndpointSample,0,len(r.samples));for _,s:=range r.samples{if s.Healthy{items=append(items,s)}};r.mu.RUnlock();if len(items)==0{return EndpointSample{},false};sort.SliceStable(items,func(i,j int)bool{if items[i].RTT!=items[j].RTT{return items[i].RTT<items[j].RTT};if items[i].Loss!=items[j].Loss{return items[i].Loss<items[j].Loss};return items[i].Endpoint<items[j].Endpoint});return items[0],true}
