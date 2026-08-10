package controlplane

import "time"

// TargetServer identifies an approved deployment destination. It is a control
// plane record; credentials are intentionally not stored in this structure.
type TargetServer struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Environment string    `json:"environment"`
	Transport   string    `json:"transport"` // agent, ssh, api
	AgentID     string    `json:"agent_id,omitempty"`
	Enabled     bool      `json:"enabled"`
	LastSeen    time.Time `json:"last_seen,omitempty"`
}

type DeploymentTarget struct {
	ServerID string `json:"server_id"`
	Path     string `json:"path"`
	Version  string `json:"version"`
}

// ValidateTarget prevents accidental deployment to an unregistered target.
func ValidateTarget(servers []TargetServer, id string) bool {
	for _, s := range servers {
		if s.ID == id && s.Enabled { return true }
	}
	return false
}
