package callcenter

import("strings";"sync")
type Handler func(CallEvent) error
type Router struct{mu sync.RWMutex;handlers map[string][]Handler}
func NewRouter()*Router{return &Router{handlers:make(map[string][]Handler)}}
func(r *Router)Register(eventType string,h Handler){if r==nil||h==nil{return};eventType=strings.TrimSpace(eventType);if eventType==""{return};r.mu.Lock();if r.handlers==nil{r.handlers=make(map[string][]Handler)};r.handlers[eventType]=append(r.handlers[eventType],h);r.mu.Unlock()}
func(r *Router)Dispatch(e CallEvent)[]error{if r==nil{return nil};e.SessionID=strings.TrimSpace(e.SessionID);e.Type=strings.TrimSpace(e.Type);if e.Type==""{return []error{nil}};r.mu.RLock();handlers:=append([]Handler(nil),r.handlers[e.Type]...);r.mu.RUnlock();errs:=make([]error,0);for _,h:=range handlers{if h==nil{continue};if err:=h(e);err!=nil{errs=append(errs,err)}};return errs}
