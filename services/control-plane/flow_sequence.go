package main

import (
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
	mu    sync.Mutex
	state map[FlowSequenceKey]FlowSequenceState
}

func NewFlowSequenceTracker() *FlowSequenceTracker {
	return &FlowSequenceTracker{state: make(map[FlowSequenceKey]FlowSequenceState)}
}

// Observe records exporter sequence continuity. A wrap-around is treated as normal.
func (t *FlowSequenceTracker) Observe(exporter string, version uint16, sequence uint32, expected uint32) (gap uint32) {
	t.mu.Lock()
	defer t.mu.Unlock()
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
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.state[FlowSequenceKey{Exporter: exporter, Version: version}]
}
