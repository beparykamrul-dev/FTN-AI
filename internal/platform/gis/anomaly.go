package gis

import "time"

type Anomaly struct {
	AssetID string `json:"asset_id"`
	Kind string `json:"kind"`
	Score float64 `json:"score"`
	Reason string `json:"reason"`
	ObservedAt time.Time `json:"observed_at"`
}

// ScoreDrift converts reconciled drift into a bounded anomaly score. This is
// deterministic infrastructure scoring; a model can be plugged in later.
func ScoreDrift(d Drift) Anomaly {
	score := 0.5
	if d.Severity == "high" { score = 0.9 }
	if d.Severity == "low" { score = 0.25 }
	return Anomaly{AssetID:d.AssetID, Kind:d.Kind, Score:score, Reason:d.Expected+" -> "+d.Observed, ObservedAt:time.Now().UTC()}
}
