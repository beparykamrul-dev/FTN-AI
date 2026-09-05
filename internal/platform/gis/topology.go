package gis
import("sort";"strings";"sync")
type TopologyLink struct{ID string `json:"id"`;FromID string `json:"from_id"`;ToID string `json:"to_id"`;Kind string `json:"kind"`;Status string `json:"status,omitempty"`}
type TopologyGraph struct{mu sync.RWMutex;links map[string]TopologyLink}
func NewTopologyGraph()*TopologyGraph{return &TopologyGraph{links:make(map[string]TopologyLink)}}
func(g *TopologyGraph)Upsert(l TopologyLink){if g==nil{return};l.ID=strings.TrimSpace(l.ID);l.FromID=strings.TrimSpace(l.FromID);l.ToID=strings.TrimSpace(l.ToID);l.Kind=strings.TrimSpace(l.Kind);l.Status=strings.ToLower(strings.TrimSpace(l.Status));if l.ID==""||l.FromID==""||l.ToID==""||l.FromID==l.ToID||len(l.ID)>256||len(l.FromID)>256||len(l.ToID)>256{return};g.mu.Lock();defer g.mu.Unlock();if g.links==nil{g.links=make(map[string]TopologyLink)};g.links[l.ID]=l}
func(g *TopologyGraph)List()[]TopologyLink{if g==nil{return nil};g.mu.RLock();defer g.mu.RUnlock();out:=make([]TopologyLink,0,len(g.links));for _,l:=range g.links{out=append(out,l)};sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID});return out}
func(g *TopologyGraph)Neighbors(id string)[]TopologyLink{if g==nil{return nil};id=strings.TrimSpace(id);if id==""{return nil};g.mu.RLock();defer g.mu.RUnlock();out:=make([]TopologyLink,0);for _,l:=range g.links{if l.FromID==id||l.ToID==id{out=append(out,l)}};sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID});return out}
