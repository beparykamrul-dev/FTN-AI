package fiber

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type FiberNode struct {
	ID, ParentID, OLTID, ONUID, ResellerID, CustomerID string
	Kind, Status string
	Latitude, Longitude float64
	Lat, Lon float64
	LastSeen time.Time
}

type FiberMap struct { Nodes map[string]FiberNode; UpdatedAt time.Time }
type RecoveryPlan struct { NodeID string; CandidateParentIDs []string; Reason string; Confidence float64; RequiresApproval bool; CreatedAt time.Time }
type RecoveryEngine struct { mu sync.RWMutex; nodes map[string]FiberNode; plans map[string]RecoveryPlan }
func NewRecoveryEngine() *RecoveryEngine { return &RecoveryEngine{nodes:map[string]FiberNode{},plans:map[string]RecoveryPlan{}} }
func (e *RecoveryEngine) Upsert(n FiberNode) error { if n.ID=="" { return fmt.Errorf("fiber node id required") }; e.mu.Lock(); e.nodes[n.ID]=n; e.mu.Unlock(); return nil }
func (e *RecoveryEngine) Snapshot(now time.Time) FiberMap { e.mu.RLock(); defer e.mu.RUnlock(); out:=map[string]FiberNode{}; for k,v:=range e.nodes { out[k]=v }; return FiberMap{Nodes:out,UpdatedAt:now} }
func (e *RecoveryEngine) Plan(nodeID string, now time.Time) (RecoveryPlan,error) { e.mu.RLock(); n,ok:=e.nodes[nodeID]; if !ok { e.mu.RUnlock(); return RecoveryPlan{},fmt.Errorf("unknown fiber node") }; var candidates []FiberNode; for _,x:=range e.nodes { if x.ID!=n.ID && x.OLTID==n.OLTID && x.Status=="up" { candidates=append(candidates,x) } }; e.mu.RUnlock(); sort.Slice(candidates,func(i,j int)bool{return candidates[i].LastSeen.After(candidates[j].LastSeen)}); ids:=make([]string,0,len(candidates)); for _,x:=range candidates { ids=append(ids,x.ID) }; p:=RecoveryPlan{NodeID:nodeID,CandidateParentIDs:ids,Reason:"topology recovery candidate set",Confidence:0.5,RequiresApproval:true,CreatedAt:now}; e.mu.Lock(); e.plans[nodeID]=p; e.mu.Unlock(); return p,nil }
