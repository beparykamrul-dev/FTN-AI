package observability

import (
	"sort"
	"strings"
	"sync"
	"time"
)

type TrafficSample struct { Interface string `json:"interface"`; SourceIP string `json:"source_ip,omitempty"`; DestIP string `json:"dest_ip,omitempty"`; Protocol string `json:"protocol,omitempty"`; App string `json:"application,omitempty"`; Bytes uint64 `json:"bytes"`; Packets uint64 `json:"packets"`; At time.Time `json:"at"` }
type TrafficStore struct { mu sync.RWMutex; samples []TrafficSample; limit int }
func NewTrafficStore(limit int) *TrafficStore { if limit<1{limit=10000};return &TrafficStore{limit:limit} }
func (s *TrafficStore) Add(v TrafficSample) { if s==nil{return};v.Interface=strings.TrimSpace(v.Interface);if v.Interface==""{return};v.Protocol=strings.ToLower(strings.TrimSpace(v.Protocol));v.App=strings.TrimSpace(v.App);if v.At.IsZero(){v.At=time.Now().UTC()};s.mu.Lock();defer s.mu.Unlock();s.samples=append(s.samples,v);if len(s.samples)>s.limit{s.samples=s.samples[len(s.samples)-s.limit:]} }
func (s *TrafficStore) List() []TrafficSample { if s==nil{return nil};s.mu.RLock();defer s.mu.RUnlock();out:=make([]TrafficSample,len(s.samples));copy(out,s.samples);sort.SliceStable(out,func(i,j int)bool{if out[i].At.Equal(out[j].At){return out[i].Interface<out[j].Interface};return out[i].At.Before(out[j].At)});return out }
