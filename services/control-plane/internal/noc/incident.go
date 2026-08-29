package noc

import "time"

type Incident struct {
	ID          string    `json:"id"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	Source      string    `json:"source"`
	DeviceID    string    `json:"device_id,omitempty"`
	AccountID   string    `json:"account_id,omitempty"`
	Summary     string    `json:"summary"`
	OccurredAt  time.Time `json:"occurred_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}
