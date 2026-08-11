package ftndns

import "sort"

type CacheNode struct { ID ProviderID; Healthy bool; LatencyMS int64; SnapshotHash string }

type CacheMeshCoordinator struct{}

func (CacheMeshCoordinator) SelectPeer(nodes []CacheNode, localHash string) (CacheNode, bool) {
	candidates:=make([]CacheNode,0,len(nodes))
	for _, n:=range nodes { if n.Healthy && n.SnapshotHash==localHash { candidates=append(candidates,n) } }
	if len(candidates)==0 { return CacheNode{},false }
	sort.SliceStable(candidates,func(i,j int) bool { if candidates[i].LatencyMS!=candidates[j].LatencyMS { return candidates[i].LatencyMS<candidates[j].LatencyMS }; return candidates[i].ID<candidates[j].ID })
	return candidates[0],true
}
