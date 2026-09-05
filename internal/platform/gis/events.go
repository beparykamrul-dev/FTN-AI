package gis
import("encoding/json";"fmt";"strings";"sync")
type Event struct{Type string `json:"type"`;Node *MapNode `json:"node,omitempty"`;Edge *MapEdge `json:"edge,omitempty"`}
type Subscriber chan []byte
type Hub struct{mu sync.RWMutex;subs map[Subscriber]struct{}}
func NewHub()*Hub{return &Hub{subs:make(map[Subscriber]struct{})}}
func(h *Hub)Subscribe(buffer int)Subscriber{if h==nil{return nil};if buffer<1{buffer=16};if buffer>4096{buffer=4096};ch:=make(Subscriber,buffer);h.mu.Lock();if h.subs==nil{h.subs=make(map[Subscriber]struct{})};h.subs[ch]=struct{}{};h.mu.Unlock();return ch}
func(h *Hub)Unsubscribe(ch Subscriber){if h==nil||ch==nil{return};h.mu.Lock();if _,ok:=h.subs[ch];ok{delete(h.subs,ch);close(ch)};h.mu.Unlock()}
func(h *Hub)Publish(e Event)error{if h==nil{return fmt.Errorf("event hub is required")};e.Type=strings.TrimSpace(e.Type);if e.Type==""||len(e.Type)>256{return fmt.Errorf("event type is invalid")};data,err:=json.Marshal(e);if err!=nil{return err};if len(data)>1<<20{return fmt.Errorf("event payload is too large")};h.mu.RLock();defer h.mu.RUnlock();for ch:=range h.subs{select{case ch<-append([]byte(nil),data...):default:}};return nil}
