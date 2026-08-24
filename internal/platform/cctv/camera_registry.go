package cctv

import (
    "fmt"
    "sync"
)

type Mode string
const (
    IP Mode = "IP_CAMERA"
    NonIP Mode = "NON_IP_CAMERA"
)

type Camera struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Mode Mode `json:"mode"`
    SiteID string `json:"site_id"`
    Address string `json:"address,omitempty"`
    Protocol string `json:"protocol,omitempty"`
    GatewayID string `json:"gateway_id,omitempty"`
    RecorderID string `json:"recorder_id,omitempty"`
    BridgeID string `json:"bridge_id,omitempty"`
    Status string `json:"status"`
}

type Registry struct { mu sync.RWMutex; cameras map[string]Camera }
func NewRegistry() *Registry { return &Registry{cameras: make(map[string]Camera)} }
func (r *Registry) Upsert(c Camera) error {
    if c.ID == "" || c.SiteID == "" { return fmt.Errorf("camera id and site_id are required") }
    if c.Mode != IP && c.Mode != NonIP { return fmt.Errorf("unsupported camera mode") }
    r.mu.Lock(); defer r.mu.Unlock(); r.cameras[c.ID] = c; return nil
}
func (r *Registry) Get(id string) (Camera, bool) { r.mu.RLock(); defer r.mu.RUnlock(); c, ok := r.cameras[id]; return c, ok }
func (r *Registry) List() []Camera { r.mu.RLock(); defer r.mu.RUnlock(); out:=make([]Camera,0,len(r.cameras)); for _,c:=range r.cameras { out=append(out,c) }; return out }
