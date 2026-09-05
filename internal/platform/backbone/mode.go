package backbone

import "strings"

type Mode string

const (
	ModeActive  Mode = "active"
	ModeStandby Mode = "standby"
)

type Backbone struct {
	PrimaryID   string `json:"primary_id"`
	SecondaryID string `json:"secondary_id"`
	Mode        Mode   `json:"mode"`
	Healthy     bool   `json:"healthy"`
}

func (b Backbone) Valid() bool {
	primary := strings.TrimSpace(b.PrimaryID)
	secondary := strings.TrimSpace(b.SecondaryID)
	return primary != "" && secondary != "" && primary != secondary && len(primary) <= 256 && len(secondary) <= 256 && (b.Mode == ModeActive || b.Mode == ModeStandby)
}

func (b Backbone) CanFailover() bool { return b.Valid() && b.Healthy && b.Mode == ModeActive }

// SwitchMode returns a proposed state only. Execution requires policy and approval.
func (b Backbone) SwitchMode() Mode {
	if b.Mode == ModeActive {
		return ModeStandby
	}
	return ModeActive
}
