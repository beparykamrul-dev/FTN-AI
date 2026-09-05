package gis

import (
	"sort"
	"strings"
	"time"
)

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

type Reconciler struct{}

func NewReconciler() *Reconciler { return &Reconciler{} }

func (r *Reconciler) Compare(expected []FiberAsset, observed []DiscoveryRecord) []Drift {
	observedByID := make(map[string]DiscoveryRecord, len(observed))
	for _, v := range observed {
		id := strings.TrimSpace(v.AssetID)
		if id != "" { observedByID[id] = v }
	}
	out := make([]Drift, 0)
	for _, e := range expected {
		id := strings.TrimSpace(e.ID)
		if id == "" { continue }
		o, ok := observedByID[id]
		if !ok {
			out = append(out, Drift{AssetID:id, Kind:string(e.Type), Expected:e.Status, Observed:"missing", Severity:"high"})
			continue
		}
		if e.Status != "" && o.Status != "" && e.Status != o.Status {
			out = append(out, Drift{AssetID:id, Kind:string(e.Type), Expected:e.Status, Observed:o.Status, Severity:"medium"})
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity { return out[i].Severity < out[j].Severity }
		if out[i].AssetID != out[j].AssetID { return out[i].AssetID < out[j].AssetID }
		return out[i].Kind < out[j].Kind
	})
	return out
}
