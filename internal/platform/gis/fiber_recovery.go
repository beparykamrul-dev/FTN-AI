package gis

import (
    "fmt"
    "sort"
)

type RecoveryCandidate struct { AssetID string `json:"asset_id"`; DistanceM float64 `json:"distance_m"`; Score float64 `json:"score"`; Reason string `json:"reason"` }

// BuildRecoveryPlan produces deterministic recovery candidates from the known
// topology. It never changes the physical topology itself.
func (m *FiberMap) BuildRecoveryPlan(failedID string) (RecoveryPlan, error) {
    if failedID == "" { return RecoveryPlan{}, fmt.Errorf("failed asset id is required") }
    assets := m.List()
    var failed *FiberAsset
    for i := range assets { if assets[i].ID == failedID { failed = &assets[i]; break } }
    if failed == nil { return RecoveryPlan{}, fmt.Errorf("fiber asset not found") }
    out := RecoveryPlan{FailedAssetID: failedID, RequiresApproval: true}
    for _, a := range assets {
        if a.ID == failedID || a.Status == "down" || a.Status == "failed" { continue }
        if a.Type != failed.Type && a.Type != FiberRoute && a.Type != FiberRouter && a.Type != FiberONU { continue }
        d := a.DistanceM; if d <= 0 { d = 1 }
        score := 1.0 / d
        out.Candidates = append(out.Candidates, RecoveryCandidate{AssetID:a.ID, DistanceM:d, Score:score, Reason:"healthy topology candidate"})
    }
    sort.Slice(out.Candidates, func(i,j int) bool { return out.Candidates[i].Score > out.Candidates[j].Score })
    return out, nil
}
