package telemetry

import("sort";"strings";"sync";"time")
type Stream struct{mu sync.RWMutex;nodes map[string]Heartbeat}
func NewStream()*Stream{return &Stream{nodes:make(map[string]Heartbeat)}}
func(s *Stream)Publish(h Heartbeat)bool{if s==nil||!h.Valid(){return false};h.NodeID=strings.TrimSpace(h.NodeID);h.ObservedAt=h.ObservedAt.UTC();s.mu.Lock();defer s.mu.Unlock();if s.nodes==nil{s.nodes=make(map[string]Heartbeat)};if old,ok:=s.nodes[h.NodeID];ok&&!h.ObservedAt.After(old.ObservedAt){return false};s.nodes[h.NodeID]=h;return true}
func(s *Stream)Snapshot(now time.Time,maxAge time.Duration)[]Heartbeat{if s==nil{return []Heartbeat{}};s.mu.RLock();out:=make([]Heartbeat,0,len(s.nodes));for _,h:=range s.nodes{if Fresh(h,now,maxAge){out=append(out,h)}};s.mu.RUnlock();sort.Slice(out,func(i,j int)bool{return out[i].NodeID<out[j].NodeID});return out}
