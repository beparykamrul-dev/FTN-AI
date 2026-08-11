package fiber

import "time"

type FaultType string

const (
	FaultCut FaultType = "cut"
	FaultDegradation FaultType = "degradation"
	FaultOpticalLoss FaultType = "optical_loss"
	FaultUnknown FaultType = "unknown"
)

type FaultEvent struct {
	ID string `json:"id"`
	PathID string `json:"path_id"`
	Type FaultType `json:"type"`
	Severity string `json:"severity"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	DistanceMeters float64 `json:"distance_meters,omitempty"`
	Summary string `json:"summary"`
	ObservedAt time.Time `json:"observed_at"`
}

// NewFaultEvent normalizes a fiber fault observation for the GIS and recovery
// layers. Location is optional because not every detector can geolocate a fault.
func NewFaultEvent(id, pathID string, kind FaultType, severity, summary string, observedAt time.Time) FaultEvent {
	return FaultEvent{ID:id, PathID:pathID, Type:kind, Severity:severity, Summary:summary, ObservedAt:observedAt.UTC()}
}
