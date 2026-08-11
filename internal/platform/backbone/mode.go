package backbone

// Mode identifies the two FTN backbone operating planes.
// Active carries normal production traffic; standby is a synchronized,
// health-checked failover plane. The control plane must require an explicit
// policy/approval before changing the active mode.
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

func (b Backbone) CanFailover() bool {
	return b.PrimaryID != "" && b.SecondaryID != "" && b.Healthy
}

// SwitchMode returns a proposed state only. Execution belongs to the
// authenticated FTN dataplane executor after policy and approval checks.
func (b Backbone) SwitchMode() Mode {
	if b.Mode == ModeActive { return ModeStandby }
	return ModeActive
}
