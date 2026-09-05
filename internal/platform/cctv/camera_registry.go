package cctv

import("fmt";"sort";"strings";"sync")
type Mode string
const(IP Mode="IP_CAMERA";NonIP Mode="NON_IP_CAMERA")
type Camera struct{ID string `json:"id"`;Name string `json:"name"`;Mode Mode `json:"mode"`;SiteID string `json:"site_id"`;Address string `json:"address,omitempty"`;Protocol string `json:"protocol,omitempty"`;GatewayID string `json:"gateway_id,omitempty"`;RecorderID string `json:"recorder_id,omitempty"`;BridgeID string `json:"bridge_id,omitempty"`;Status string `json:"status"`}
type Registry struct{mu sync.RWMutex;cameras map[string]Camera}
func NewRegistry()*Registry{return &Registry{cameras:make(map[string]Camera)}}
func(r *Registry)Upsert(c Camera)error{if r==nil{return fmt.Errorf("camera registry is required")};c.ID=strings.TrimSpace(c.ID);c.SiteID=strings.TrimSpace(c.SiteID);if c.ID==""||c.SiteID==""{return fmt.Errorf("camera id and site_id are required")};if c.Mode!=IP&&c.Mode!=NonIP{return fmt.Errorf("unsupported camera mode")};r.mu.Lock();defer r.mu.Unlock();if r.cameras==nil{r.cameras=make(map[string]Camera)};r.cameras[c.ID]=c;return nil}
func(r *Registry)Get(id string)(Camera,bool){if r==nil{return Camera{},false};id=strings.TrimSpace(id);r.mu.RLock();defer r.mu.RUnlock();c,ok:=r.cameras[id];return c,ok}
func(r *Registry)List()[]Camera{if r==nil{return nil};r.mu.RLock();out:=make([]Camera,0,len(r.cameras));for _,c:=range r.cameras{out=append(out,c)};r.mu.RUnlock();sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID});return out}
