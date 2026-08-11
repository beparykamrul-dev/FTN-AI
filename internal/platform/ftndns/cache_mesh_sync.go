package ftndns

import "sort"

type CacheSyncPlan struct { Source ProviderID; Targets []ProviderID; SnapshotHash string }

// BuildCacheSyncPlan only plans synchronization between healthy peers that
// share the same normalized snapshot. Execution remains outside the planner.
func BuildCacheSyncPlan(localHash string, local ProviderID, nodes []CacheNode) CacheSyncPlan {
	candidates:=make([]CacheNode,0,len(nodes))
	for _, n:=range nodes { if n.ID!=local && n.Healthy && n.SnapshotHash==localHash { candidates=append(candidates,n) } }
	sort.SliceStable(candidates,func(i,j int) bool { if candidates[i].LatencyMS!=candidates[j].LatencyMS { return candidates[i].LatencyMS<candidates[j].LatencyMS }; return candidates[i].ID<candidates[j].ID })
	plan:=CacheSyncPlan{Source:local,SnapshotHash:localHash}; for _,n:=range candidates { plan.Targets=append(plan.Targets,n.ID) }; return plan
}
