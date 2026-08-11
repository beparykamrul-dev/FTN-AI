package fiber

import "time"

type FiberState string

const (
	FiberHealthy  FiberState = "healthy"
	FiberDegraded FiberState = "degraded"
	FiberDown     FiberState = "down"
	FiberUnknown  FiberState = "unknown"
)

// Segment is the canonical FTN representation of a monitored fiber segment.
type Segment struct {
	ID             string     `json:"id"`
	FromNode       string     `json:"from_node"`
	ToNode         string     `json:"to_node"`
	DistanceMeters float64    `json:"distance_meters"`
	State          FiberState `json:"state"`
	RxPowerDBm     *float64   `json:"rx_power_dbm,omitempty"`
	TxPowerDBm     *float64   `json:"tx_power_dbm,omitempty"`
	LossDB         *float64   `json:"loss_db,omitempty"`
	LastSeen       time.Time  `json:"last_seen"`
}

type FiberAlert struct {
	SegmentID string    `json:"segment_id"`
	Severity  string    `json:"severity"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// Evaluate classifies an observed segment. It produces state/alert data only;
// repair or traffic changes are handled by the approved recovery workflow.
func Evaluate(s Segment, now time.Time) (Segment, *FiberAlert) {
	s.LastSeen = s.LastSeen.UTC()
	if s.State == FiberDown || s.State == FiberDegraded {
		return s, &FiberAlert{SegmentID:s.ID, Severity:string(s.State), Reason:"fiber segment requires investigation", CreatedAt:now.UTC()}
	}
	return s, nil
}
