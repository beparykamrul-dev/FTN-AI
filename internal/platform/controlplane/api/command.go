package api

import "time"

type CommandKind string

const (
	CommandObserve CommandKind = "observe"
	CommandPlan    CommandKind = "plan"
	CommandExecute CommandKind = "execute"
)

type Command struct {
	ID         string      `json:"id"`
	Kind       CommandKind `json:"kind"`
	Target     string      `json:"target"`
	Action     string      `json:"action"`
	RequestedBy string     `json:"requested_by"`
	CreatedAt  time.Time   `json:"created_at"`
	Approved   bool        `json:"approved"`
}

// CanExecute enforces FTN's approval-first control boundary.
func (c Command) CanExecute() bool {
	return c.ID != "" && c.Target != "" && c.Action != "" && c.Approved && c.Kind == CommandExecute
}
