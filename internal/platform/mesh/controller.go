package mesh

import (
    "sort"
    "sync"
    "time"
)

type Node struct { ID, Address string; Online bool; LastSeen time.Time }
type Link struct { ID, From, To string; Metric uint32; Up bool; LastSeen time.Time }
type Snapshot struct { Nodes []Node `json:"nodes"`; Links []Link `json:"links"`; GeneratedAt time.Time `json:"generated_at"` }

type Controller struct { mu sync.RWMutex; nodes map[string]Node; links map[string]Link }
func NewController() *Controller { return &Controller{nodes: map[string]Node{}, links: map[string]Link{}} }
func (c *Controller) UpsertNode(n Node) { c.mu.Lock(); defer c.mu.Unlock(); c.nodes[n.ID]=n }
func (c *Controller) UpsertLink(l Link) { c.mu.Lock(); defer c.mu.Unlock(); c.links[l.ID]=l }
func (c *Controller) Snapshot() Snapshot { c.mu.RLock(); defer c.mu.RUnlock(); n:=make([]Node,0,len(c.nodes)); for _,v:=range c.nodes { n=append(n,v) }; l:=make([]Link,0,len(c.links)); for _,v:=range c.links { l=append(l,v) }; sort.Slice(n,func(i,j int)bool{return n[i].ID<n[j].ID}); sort.Slice(l,func(i,j int)bool{return l[i].ID<l[j].ID}); return Snapshot{Nodes:n,Links:l,GeneratedAt:time.Now().UTC()} }
