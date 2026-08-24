package media

import (
 "fmt"
 "sync"
 "time"
)

type RuntimeState struct { ItemID string `json:"item_id"`; State string `json:"state"`; LastCheck time.Time `json:"last_check"`; Error string `json:"error,omitempty"` }
type ServiceRuntime struct { mu sync.RWMutex; states map[string]RuntimeState }
func NewServiceRuntime()*ServiceRuntime{return &ServiceRuntime{states:map[string]RuntimeState{}}}
func(r *ServiceRuntime) Update(id,state string,err error){r.mu.Lock();defer r.mu.Unlock();v:=RuntimeState{ItemID:id,State:state,LastCheck:time.Now()};if err!=nil{v.Error=err.Error()};r.states[id]=v}
func(r *ServiceRuntime) Check(id string)(RuntimeState,error){r.mu.RLock();defer r.mu.RUnlock();v,ok:=r.states[id];if !ok{return RuntimeState{},fmt.Errorf("media runtime not found")};return v,nil}
func(r *ServiceRuntime) List()[]RuntimeState{r.mu.RLock();defer r.mu.RUnlock();o:=make([]RuntimeState,0,len(r.states));for _,v:=range r.states{o=append(o,v)};return o}
