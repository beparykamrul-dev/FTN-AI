package mesh

import (
	"sync"
	"time"
)

type Keepalive struct {
	NodeID    string    `json:"node_id"`
	SentAt    time.Time `json:"sent_at"`
	Sequence  uint64    `json:"sequence"`
}

type NodeState struct {
	NodeID   string    `json:"node_id"`
	LastSeen time.Time `json:"last_seen"`
	Healthy  bool      `json:"healthy"`
	Sequence uint64    `json:"sequence"`
}

type StateStore struct {
	mu    sync.RWMutex
	nodes map[string]NodeState
}

func NewStateStore() *StateStore { return &StateStore{nodes: make(map[string]NodeState)} }

func (s *StateStore) Observe(k Keepalive) NodeState {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := NodeState{NodeID: k.NodeID, LastSeen: k.SentAt.UTC(), Healthy: true, Sequence: k.Sequence}
	s.nodes[k.NodeID] = state
	return state
}

func (s *StateStore) Healthy(nodeID string, now time.Time, timeout time.Duration) bool {
	s.mu.RLock()
	state, ok := s.nodes[nodeID]
	s.mu.RUnlock()
	return ok && state.Healthy && !state.LastSeen.IsZero() && now.Sub(state.LastSeen) <= timeout
}
