package agent
import("sort";"strings";"sync")
type Adapter interface{Kind() DeviceKind;Capabilities()[]string}
type AdapterRegistry struct{mu sync.RWMutex;adapters map[DeviceKind]Adapter}
func NewAdapterRegistry()*AdapterRegistry{return &AdapterRegistry{adapters:make(map[DeviceKind]Adapter)}}
func(r *AdapterRegistry)Register(a Adapter){if r==nil||a==nil{return};kind:=DeviceKind(strings.TrimSpace(string(a.Kind())));if kind==""{return};r.mu.Lock();defer r.mu.Unlock();if r.adapters==nil{r.adapters=make(map[DeviceKind]Adapter)};r.adapters[kind]=a}
func(r *AdapterRegistry)Get(kind DeviceKind)(Adapter,bool){if r==nil{return nil,false};kind=DeviceKind(strings.TrimSpace(string(kind)));if kind==""{return nil,false};r.mu.RLock();defer r.mu.RUnlock();a,ok:=r.adapters[kind];return a,ok}
func(r *AdapterRegistry)Kinds()[]DeviceKind{if r==nil{return nil};r.mu.RLock();defer r.mu.RUnlock();out:=make([]DeviceKind,0,len(r.adapters));for k:=range r.adapters{out=append(out,k)};sort.Slice(out,func(i,j int)bool{return out[i]<out[j]});return out}
