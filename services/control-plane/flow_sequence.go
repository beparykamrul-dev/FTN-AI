package main

import (
	"strings"
	"sync"
)

type FlowSequenceKey struct {
	Exporter string
	Version  uint16
}

type FlowSequenceState struct {
	Initialized bool
	Last       uint32
	Gaps       uint64
	Packets    uint64
}

type FlowSequenceTracker struct {
	mu    sync.RWMutex
	state map[FlowSequenceKey]FlowSequenceState
}

func NewFlowSequenceTracker() *FlowSequenceTracker {
	return &FlowSequenceTracker{state: make(map[FlowSequenceKey]FlowSequenceState)}
}

// Observe records exporter sequence continuity. A wrap-around is treated as normal.
func (t *FlowSequenceTracker) Observe(exporter string, version uint16, sequence uint32, expected uint32) (gap uint32) {
	if t == nil { return 0 }
	exporter = strings.TrimSpace(exporter)
	if exporter == "" { return 0 }
	if expected == 0 { expected = 1 }
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == nil { t.state = make(map[FlowSequenceKey]FlowSequenceState) }
	key := FlowSequenceKey{Exporter: exporter, Version: version}
	s := t.state[key]
	if !s.Initialized {
		s.Initialized = true
		s.Last = sequence
		s.Packets++
		t.state[key] = s
		return 0
	}
	actual := sequence - s.Last
	if actual > expected {
		gap = actual - expected
		s.Gaps += uint64(gap)
	}
	s.Last = sequence
	s.Packets++
	t.state[key] = s
	return gap
}

func (t *FlowSequenceTracker) State(exporter string, version uint16) FlowSequenceState {
	if t == nil { return FlowSequenceState{} }
	exporter = strings.TrimSpace(exporter)
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.state[FlowSequenceKey{Exporter: exporter, Version: version}]
}
