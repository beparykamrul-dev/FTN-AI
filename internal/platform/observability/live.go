package observability

import (
	"encoding/json"
	"sync"
)

type LiveEvent struct { Type string `json:"type"`; Data any `json:"data,omitempty"` }
type LiveHub struct { mu sync.RWMutex; subs map[chan []byte]struct{} }
func NewLiveHub() *LiveHub { return &LiveHub{subs: make(map[chan []byte]struct{})} }
func (h *LiveHub) Subscribe(buffer int) chan []byte { if h==nil{return nil};if buffer<1{buffer=32};ch:=make(chan []byte,buffer);h.mu.Lock();if h.subs==nil{h.subs=make(map[chan []byte]struct{})};h.subs[ch]=struct{}{};h.mu.Unlock();return ch }
func (h *LiveHub) Unsubscribe(ch chan []byte) { if h==nil||ch==nil{return};h.mu.Lock();if _,ok:=h.subs[ch];ok{delete(h.subs,ch);close(ch)};h.mu.Unlock() }
func (h *LiveHub) Publish(e LiveEvent) error { data,err:=json.Marshal(e);if err!=nil{return err};if h==nil{return nil};h.mu.RLock();defer h.mu.RUnlock();for ch:=range h.subs{select{case ch<-append([]byte(nil),data...):default:}};return nil }
