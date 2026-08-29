package audit

import "time"

type Event struct {
	ID         string    `json:"id"`
	RequestID  string    `json:"request_id"`
	ActorID    string    `json:"actor_id,omitempty"`
	ActorRole  string    `json:"actor_role,omitempty"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID string    `json:"resource_id,omitempty"`
	Outcome    string    `json:"outcome"`
	Timestamp  time.Time `json:"timestamp"`
}

const (
	OutcomeSuccess = "success"
	OutcomeDenied  = "denied"
	OutcomeFailed  = "failed"
)
