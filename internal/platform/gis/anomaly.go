package gis

import (
	"strings"
	"time"
)

type Anomaly struct {
	AssetID     string    `json:"asset_id"`
	Kind        string    `json:"kind"`
	Score       float64   `json:"score"`
	Reason      string    `json:"reason"`
	ObservedAt  time.Time `json:"observed_at"`
}

// ScoreDrift converts reconciled drift into a bounded anomaly score. This is
// deterministic infrastructure scoring; a model can be plugged in later.
func ScoreDrift(d Drift) Anomaly {
	score := 0.5
	switch strings.ToLower(strings.TrimSpace(d.Severity)) {
	case "high":
		score = 0.9
	case "low":
		score = 0.25
	}
	return Anomaly{AssetID: strings.TrimSpace(d.AssetID), Kind: strings.TrimSpace(d.Kind), Score: score, Reason: strings.TrimSpace(d.Expected) + " -> " + strings.TrimSpace(d.Observed), ObservedAt: time.Now().UTC()}
}
