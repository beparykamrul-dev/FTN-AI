package telemetry

import("encoding/json";"strings";"sync")
type WebSocketHub struct{mu sync.RWMutex;sessions map[string]chan []byte}
func NewWebSocketHub()*WebSocketHub{return &WebSocketHub{sessions:make(map[string]chan []byte)}}
func(h *WebSocketHub)Register(sessionID string,buffer int)chan []byte{if h==nil{return nil};sessionID=strings.TrimSpace(sessionID);if sessionID==""{return nil};if buffer<1{buffer=1};ch:=make(chan []byte,buffer);h.mu.Lock();defer h.mu.Unlock();if h.sessions==nil{h.sessions=make(map[string]chan []byte)};if old,ok:=h.sessions[sessionID];ok{close(old)};h.sessions[sessionID]=ch;return ch}
func(h *WebSocketHub)Unregister(sessionID string){if h==nil{return};sessionID=strings.TrimSpace(sessionID);if sessionID==""{return};h.mu.Lock();defer h.mu.Unlock();if ch,ok:=h.sessions[sessionID];ok{close(ch);delete(h.sessions,sessionID)}}
func(h *WebSocketHub)BroadcastHeartbeat(hb Heartbeat)int{if h==nil||!hb.Valid(){return 0};payload,err:=json.Marshal(hb);if err!=nil{return 0};h.mu.RLock();defer h.mu.RUnlock();sent:=0;for _,ch:=range h.sessions{select{case ch<-payload:sent++;default:}};return sent}
