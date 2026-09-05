package agent

import "time"

type Heartbeat struct {
	AgentID    string    `json:"agent_id"`
	HostID     string    `json:"host_id"`
	Version    string    `json:"version"`
	Status     string    `json:"status"`
	ObservedAt time.Time `json:"observed_at"`
}

type AgentState struct {
	AgentID  string    `json:"agent_id"`
	HostID   string    `json:"host_id"`
	Version  string    `json:"version"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

// Registry accepts authenticated agent observations. Authentication and
// transport security belong to the FTN mTLS/WebSocket gateway.
type Registry struct {
	states map[string]AgentState
}

func NewRegistry() *Registry {
	return &Registry{states: make(map[string]AgentState)}
}

func (r *Registry) Upsert(h Heartbeat) AgentState {
	s := AgentState{
		AgentID:  h.AgentID,
		HostID:   h.HostID,
		Version:  h.Version,
		Status:   h.Status,
		LastSeen: h.ObservedAt.UTC(),
	}
	r.states[h.AgentID] = s
	return s
}

func (r *Registry) Get(agentID string) (AgentState, bool) {
	s, ok := r.states[agentID]
	return s, ok
}
