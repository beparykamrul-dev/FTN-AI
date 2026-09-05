package dns
import("sort";"strings";"sync";"time")
type FailoverNode struct{ID string `json:"id"`;Address string `json:"address"`;Priority uint32 `json:"priority"`;Healthy bool `json:"healthy"`;LastCheck time.Time `json:"last_check"`}
type FailoverPlan struct{Zone string `json:"zone"`;Primary string `json:"primary"`;Backups []string `json:"backups"`}
type FailoverController struct{mu sync.RWMutex;nodes map[string]FailoverNode}
func NewFailoverController()*FailoverController{return &FailoverController{nodes:make(map[string]FailoverNode)}}
func(c *FailoverController)Upsert(node FailoverNode){if c==nil{return};node.ID=strings.TrimSpace(node.ID);node.Address=strings.TrimSpace(node.Address);if node.ID==""||node.Address==""||len(node.ID)>256||len(node.Address)>256{return};c.mu.Lock();defer c.mu.Unlock();if c.nodes==nil{c.nodes=make(map[string]FailoverNode)};c.nodes[node.ID]=node}
func(c *FailoverController)Plan(zone string)FailoverPlan{zone=strings.TrimSpace(zone);if c==nil{return FailoverPlan{Zone:zone}};c.mu.RLock();defer c.mu.RUnlock();nodes:=make([]FailoverNode,0,len(c.nodes));for _,n:=range c.nodes{n.ID=strings.TrimSpace(n.ID);n.Address=strings.TrimSpace(n.Address);if n.Healthy&&n.ID!=""&&n.Address!=""{nodes=append(nodes,n)}};sort.Slice(nodes,func(i,j int)bool{if nodes[i].Priority!=nodes[j].Priority{return nodes[i].Priority<nodes[j].Priority};return nodes[i].ID<nodes[j].ID});plan:=FailoverPlan{Zone:zone};if len(nodes)==0{return plan};plan.Primary=nodes[0].ID;for _,n:=range nodes[1:]{plan.Backups=append(plan.Backups,n.ID)};return plan}
