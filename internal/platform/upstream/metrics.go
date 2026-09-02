package upstream

import (
	"sync"
	"time"
)

// Metrics is an in-process, dependency-free snapshot for exposing upstream
// state to the FTN monitoring/control plane. It intentionally contains no
// credentials or payload data.
type Metrics struct {
	mu sync.RWMutex
	m  map[string]Snapshot
}

type Snapshot struct {
	Healthy       bool
	Latency       time.Duration
	PacketLossPct float64
	PrefixCount   uint32
	PrefixLimit   uint32
	LastChecked   time.Time
	State         string
}

func NewMetrics() *Metrics { return &Metrics{m: make(map[string]Snapshot)} }

func (m *Metrics) Set(name string, s Snapshot) {
	if m == nil || name == "" { return }
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[name] = s
}

func (m *Metrics) Get(name string) (Snapshot, bool) {
	if m == nil { return Snapshot{}, false }
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.m[name]
	return s, ok
}

func (m *Metrics) All() map[string]Snapshot {
	out := make(map[string]Snapshot)
	if m == nil { return out }
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.m { out[k] = v }
	return out
}
