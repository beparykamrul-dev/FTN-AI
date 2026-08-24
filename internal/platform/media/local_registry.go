package media

import "sync"

type Item struct { ID, Title, Type, Region, Source, OriginNode, CachePolicy, DRM, Status string }
type Registry struct { mu sync.RWMutex; items map[string]Item }
func NewRegistry()*Registry{return &Registry{items:map[string]Item{}}}
func(r *Registry) Upsert(v Item){r.mu.Lock();defer r.mu.Unlock();r.items[v.ID]=v}
func(r *Registry) Get(id string)(Item,bool){r.mu.RLock();defer r.mu.RUnlock();v,ok:=r.items[id];return v,ok}
func(r *Registry) List()[]Item{r.mu.RLock();defer r.mu.RUnlock();o:=make([]Item,0,len(r.items));for _,v:=range r.items{o=append(o,v)};return o}
