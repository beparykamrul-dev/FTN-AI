package mesh

import (
	"sort"
	"strings"
	"sync"
	"time"
)

type RouteState struct { Destination string `json:"destination"`; NextHop string `json:"next_hop"`; Metric uint32 `json:"metric"`; UpdatedAt time.Time `json:"updated_at"` }
type RouteTable struct { mu sync.RWMutex; routes map[string]RouteState }
func NewRouteTable() *RouteTable { return &RouteTable{routes: make(map[string]RouteState)} }
func (r *RouteTable) Upsert(route RouteState) { if r == nil { return }; route.Destination=strings.TrimSpace(route.Destination); route.NextHop=strings.TrimSpace(route.NextHop); if route.Destination=="" { return }; if route.Destination==route.NextHop { return }; if route.UpdatedAt.IsZero(){route.UpdatedAt=time.Now().UTC()}; r.mu.Lock(); defer r.mu.Unlock(); r.routes[route.Destination]=route }
func (r *RouteTable) Get(destination string)(RouteState,bool){ if r==nil{return RouteState{},false}; r.mu.RLock(); defer r.mu.RUnlock(); v,ok:=r.routes[strings.TrimSpace(destination)]; return v,ok }
func (r *RouteTable) Snapshot()[]RouteState{ if r==nil{return []RouteState{}}; r.mu.RLock(); defer r.mu.RUnlock(); out:=make([]RouteState,0,len(r.routes)); for _,v:=range r.routes{out=append(out,v)}; sort.SliceStable(out,func(i,j int)bool{if out[i].Destination!=out[j].Destination{return out[i].Destination<out[j].Destination};return out[i].NextHop<out[j].NextHop}); return out }
