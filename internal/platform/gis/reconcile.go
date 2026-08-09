package gis

import "time"

type DiscoveryRecord struct {
	AssetID string `json:"asset_id"`
	Kind string `json:"kind"`
	Name string `json:"name"`
	Status string `json:"status"`
	ObservedAt time.Time `json:"observed_at"`
}

type Drift struct {
	AssetID string `json:"asset_id"`
	Kind string `json:"kind"`
	Expected string `json:"expected"`
	Observed string `json:"observed"`
	Severity string `json:"severity"`
}

// Reconciler compares modeled assets with authorized discovery observations.
// It reports drift only; it never changes network state.
type Reconciler struct{}

func NewReconciler() *Reconciler { return &Reconciler{} }

func (r *Reconciler) Compare(expected []FiberAsset, observed []DiscoveryRecord) []Drift {
	m := make(map[string]DiscoveryRecord, len(observed))
	for _, v := range observed { m[v.AssetID] = v }
	out := make([]Drift, 0)
	for _, e := range expected {
		o, ok := m[e.ID]
		if !ok { out = append(out, Drift{AssetID:e.ID, Kind:string(e.Type), Expected:e.Status, Observed:"missing", Severity:"high"}); continue }
		if e.Status != "" && o.Status != "" && e.Status != o.Status { out = append(out, Drift{AssetID:e.ID, Kind:string(e.Type), Expected:e.Status, Observed:o.Status, Severity:"medium"}) }
	}
	return out
}
