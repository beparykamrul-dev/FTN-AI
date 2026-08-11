package telemetry

import (
	"sync"
	"time"
)

// Stream keeps the latest validated heartbeat per FTN node. A transport such
// as WebSocket can publish these snapshots without coupling transport code to
// fleet-balancing logic.
type Stream struct {
	mu    sync.RWMutex
	nodes map[string]Heartbeat
}

func NewStream() *Stream { return &Stream{nodes: make(map[string]Heartbeat)} }

func (s *Stream) Publish(h Heartbeat) bool {
	if !h.Valid() { return false }
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.nodes[h.NodeID]; ok && h.ObservedAt.Before(old.ObservedAt) { return false }
	s.nodes[h.NodeID] = h
	return true
}

func (s *Stream) Snapshot(now time.Time, maxAge time.Duration) []Heartbeat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Heartbeat, 0, len(s.nodes))
	for _, h := range s.nodes {
		if Fresh(h, now, maxAge) { out = append(out, h) }
	}
	return out
}
