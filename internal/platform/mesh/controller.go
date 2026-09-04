package mesh

import (
    "sort"
    "sync"
    "time"
)

type Node struct {
    ID string `json:"id"`
    Name string `json:"name,omitempty"`
    Address string `json:"address,omitempty"`
    Role string `json:"role,omitempty"`
    Region string `json:"region,omitempty"`
    Endpoint string `json:"endpoint,omitempty"`
    Scope Scope `json:"scope,omitempty"`
    Enabled bool `json:"enabled"`
    Online bool `json:"online"`
    LastSeen time.Time `json:"last_seen,omitempty"`
}

type Link struct {
    ID string `json:"id,omitempty"`
    From string `json:"from"`
    To string `json:"to"`
    Metric uint32 `json:"metric"`
    State LinkState `json:"state,omitempty"`
    Up bool `json:"up,omitempty"`
    Healthy bool `json:"healthy,omitempty"`
    LatencyMS float64 `json:"latency_ms,omitempty"`
    RTTMillis float64 `json:"rttMillis,omitempty"`
    LossPercent float64 `json:"loss_percent,omitempty"`
    Loss float64 `json:"loss,omitempty"`
    JitterMs float64 `json:"jitterMs,omitempty"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
    LastSeen time.Time `json:"last_seen,omitempty"`
}

type Snapshot struct { Nodes []Node `json:"nodes"`; Links []Link `json:"links"`; GeneratedAt time.Time `json:"generated_at"` }

type Controller struct { mu sync.RWMutex; nodes map[string]Node; links map[string]Link }
func NewController() *Controller { return &Controller{nodes: map[string]Node{}, links: map[string]Link{}} }
func (c *Controller) UpsertNode(n Node) { c.mu.Lock(); defer c.mu.Unlock(); c.nodes[n.ID]=n }
func (c *Controller) UpsertLink(l Link) { c.mu.Lock(); defer c.mu.Unlock(); c.links[l.ID]=l }
func (c *Controller) Snapshot() Snapshot { c.mu.RLock(); defer c.mu.RUnlock(); n:=make([]Node,0,len(c.nodes)); for _,v:=range c.nodes { n=append(n,v) }; l:=make([]Link,0,len(c.links)); for _,v:=range c.links { l=append(l,v) }; sort.Slice(n,func(i,j int)bool{return n[i].ID<n[j].ID}); sort.Slice(l,func(i,j int)bool{return l[i].ID<l[j].ID}); return Snapshot{Nodes:n,Links:l,GeneratedAt:time.Now().UTC()} }
