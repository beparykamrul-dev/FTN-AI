package mesh

import (
	"sort"
	"sync"
	"time"
)

type LinkState string

const (
	LinkUp       LinkState = "up"
	LinkDown     LinkState = "down"
	LinkDegraded LinkState = "degraded"
)

type Link struct {
	ID          string    `json:"id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	State       LinkState `json:"state"`
	LatencyMS   float64   `json:"latency_ms"`
	LossPercent float64   `json:"loss_percent"`
	Metric      uint32    `json:"metric"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LinkStateStore struct {
	mu    sync.RWMutex
	links map[string]Link
}

func NewLinkStateStore() *LinkStateStore { return &LinkStateStore{links: make(map[string]Link)} }

func (s *LinkStateStore) Upsert(l Link) {
	if l.ID == "" || l.From == "" || l.To == "" || l.From == l.To { return }
	if l.UpdatedAt.IsZero() { l.UpdatedAt = time.Now().UTC() }
	if l.Metric == 0 { l.Metric = healthMetric(l.LatencyMS, l.LossPercent) }
	s.mu.Lock()
	s.links[l.ID] = l
	s.mu.Unlock()
}

func healthMetric(latencyMS, lossPercent float64) uint32 {
	if latencyMS < 0 { latencyMS = 0 }
	if lossPercent < 0 { lossPercent = 0 }
	v := 100.0 + latencyMS + lossPercent*20
	if v > 4294967295 { return ^uint32(0) }
	return uint32(v)
}

func (s *LinkStateStore) Get(id string) (Link, bool) {
	s.mu.RLock(); defer s.mu.RUnlock()
	l, ok := s.links[id]
	return l, ok
}

func (s *LinkStateStore) Snapshot() []Link {
	s.mu.RLock()
	out := make([]Link, 0, len(s.links))
	for _, l := range s.links { out = append(out, l) }
	s.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *LinkStateStore) Healthy(id string) bool {
	l, ok := s.Get(id)
	return ok && l.State == LinkUp
}
