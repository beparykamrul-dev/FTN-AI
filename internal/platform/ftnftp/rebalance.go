package ftnftp

import "sort"

// NodeCapacity describes runtime capacity used by the placement engine.
type NodeCapacity struct {
    ID string `json:"id"`
    Healthy bool `json:"healthy"`
    FreeBytes int64 `json:"free_bytes"`
    UsedBytes int64 `json:"used_bytes"`
    Weight float64 `json:"weight"`
}

type Placement struct {
    ObjectID string `json:"object_id"`
    Primary string `json:"primary"`
    Replicas []string `json:"replicas"`
}

// ChoosePlacement prefers healthy nodes with the most free capacity and keeps
// replicas on distinct nodes. It returns a plan only; the caller performs the
// actual transfer/commit operation.
func ChoosePlacement(objectID string, nodes []NodeCapacity, replicaCount int) Placement {
    if replicaCount < 0 { replicaCount = 0 }
    candidates := make([]NodeCapacity, 0, len(nodes))
    for _, n := range nodes { if n.Healthy && n.FreeBytes > 0 { candidates = append(candidates, n) } }
    sort.Slice(candidates, func(i,j int) bool {
        if candidates[i].Weight == candidates[j].Weight { return candidates[i].FreeBytes > candidates[j].FreeBytes }
        return candidates[i].Weight > candidates[j].Weight
    })
    p := Placement{ObjectID: objectID}
    if len(candidates) == 0 { return p }
    p.Primary = candidates[0].ID
    for i:=1; i<len(candidates) && len(p.Replicas)<replicaCount; i++ { p.Replicas = append(p.Replicas, candidates[i].ID) }
    return p
}
