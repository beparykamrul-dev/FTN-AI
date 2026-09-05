package agent

import (
	"sort"
	"sync"
	"time"
)

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
	mu     sync.RWMutex
	states map[string]AgentState
}

func NewRegistry() *Registry {
	return &Registry{states: make(map[string]AgentState)}
}

func (r *Registry) Upsert(h Heartbeat) AgentState {
	if r == nil || h.AgentID == "" {
		return AgentState{}
	}
	s := AgentState{
		AgentID:  h.AgentID,
		HostID:   h.HostID,
		Version:  h.Version,
		Status:   h.Status,
		LastSeen: h.ObservedAt.UTC(),
	}
	r.mu.Lock()
	r.states[h.AgentID] = s
	r.mu.Unlock()
	return s
}

func (r *Registry) Get(agentID string) (AgentState, bool) {
	if r == nil {
		return AgentState{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.states[agentID]
	return s, ok
}

func (r *Registry) List() []AgentState {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]AgentState, 0, len(r.states))
	for _, s := range r.states {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].AgentID < out[j].AgentID })
	return out
}
