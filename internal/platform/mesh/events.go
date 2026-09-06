package mesh

import("strings";"sync";"time")
type EventType string
const(LinkStateChanged EventType="mesh.link_state_changed";TopologyChanged EventType="mesh.topology_changed";Heartbeat EventType="mesh.heartbeat")
type Event struct{Type EventType `json:"type"`;NodeID string `json:"node_id,omitempty"`;Link Link `json:"link,omitempty"`;At time.Time `json:"at"`}
type EventBus struct{mu sync.RWMutex;subs map[chan Event]struct{}}
func NewEventBus()*EventBus{return &EventBus{subs:make(map[chan Event]struct{})}}
func(b *EventBus)Subscribe(buffer int)(<-chan Event,func()){if b==nil{return nil,func(){}};if buffer<1{buffer=16};ch:=make(chan Event,buffer);b.mu.Lock();if b.subs==nil{b.subs=make(map[chan Event]struct{})};b.subs[ch]=struct{}{};b.mu.Unlock();return ch,func(){b.mu.Lock();if _,ok:=b.subs[ch];ok{delete(b.subs,ch);close(ch)};b.mu.Unlock()}}
func(b *EventBus)Publish(e Event){if b==nil{return};e.NodeID=strings.TrimSpace(e.NodeID);if e.At.IsZero(){e.At=time.Now().UTC()}else{e.At=e.At.UTC()};e.Type=EventType(strings.TrimSpace(string(e.Type)));if e.Type==""{return};b.mu.RLock();defer b.mu.RUnlock();for ch:=range b.subs{select{case ch<-e:default:}}}
