package recovery

import "time"

// Incident is a normalized observation from any FTN subsystem.
type Incident struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	Severity  string    `json:"severity"`
	TargetID  string    `json:"target_id"`
	Summary   string    `json:"summary"`
	ObservedAt time.Time `json:"observed_at"`
}

// Action is a proposed recovery operation. Proposals are inert until an
// authorized policy/approval workflow explicitly accepts them.
type Action struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	TargetID    string `json:"target_id"`
	Description string `json:"description"`
	Risk        string `json:"risk"`
}

type RecoveryPlan struct {
	IncidentID string   `json:"incident_id"`
	Actions    []Action `json:"actions"`
	CreatedAt  time.Time `json:"created_at"`
	RequiresApproval bool `json:"requires_approval"`
}

// Planner creates conservative, explainable proposals. It never executes a
// network/server change and never treats device identifiers as credentials.
func Planner(i Incident) RecoveryPlan {
	p := RecoveryPlan{IncidentID:i.ID, CreatedAt:time.Now().UTC(), RequiresApproval:true}
	if i.TargetID != "" {
		p.Actions = append(p.Actions, Action{
			ID: i.ID+":health-check",
			Type: "health_check",
			TargetID: i.TargetID,
			Description: "re-check target health and dependent links before recovery",
			Risk: "low",
		})
	}
	return p
}
