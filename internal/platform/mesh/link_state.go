package mesh

import (
	"sort"
	"strings"
	"sync"
	"time"
)

type LinkState string

const (
	LinkUp       LinkState = "up"
	LinkDown     LinkState = "down"
	LinkDegraded LinkState = "degraded"
)

type LinkStateStore struct {
	mu    sync.RWMutex
	links map[string]Link
}

func NewLinkStateStore() *LinkStateStore { return &LinkStateStore{links: make(map[string]Link)} }

func (s *LinkStateStore) Upsert(l Link) {
	if s == nil {
		return
	}
	l.ID = strings.TrimSpace(l.ID)
	l.From = strings.TrimSpace(l.From)
	l.To = strings.TrimSpace(l.To)
	if l.ID == "" || l.From == "" || l.To == "" || l.From == l.To {
		return
	}
	if l.UpdatedAt.IsZero() {
		l.UpdatedAt = time.Now().UTC()
	}
	if l.Metric == 0 {
		l.Metric = healthMetric(l.LatencyMS, l.LossPercent)
	}
	s.mu.Lock()
	s.links[l.ID] = l
	s.mu.Unlock()
}

func healthMetric(latencyMS, lossPercent float64) uint32 {
	if latencyMS < 0 {
		latencyMS = 0
	}
	if lossPercent < 0 {
		lossPercent = 0
	}
	v := 100.0 + latencyMS + lossPercent*20
	if v > 4294967295 {
		return ^uint32(0)
	}
	return uint32(v)
}

func (s *LinkStateStore) Get(id string) (Link, bool) {
	if s == nil {
		return Link{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.links[strings.TrimSpace(id)]
	return l, ok
}

func (s *LinkStateStore) Snapshot() []Link {
	if s == nil {
		return []Link{}
	}
	s.mu.RLock()
	out := make([]Link, 0, len(s.links))
	for _, l := range s.links {
		out = append(out, l)
	}
	s.mu.RUnlock()
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].ID != out[j].ID {
			return out[i].ID < out[j].ID
		}
		return out[i].UpdatedAt.Before(out[j].UpdatedAt)
	})
	return out
}

func (s *LinkStateStore) Healthy(id string) bool {
	l, ok := s.Get(id)
	return ok && l.State == LinkUp
}

// HealthyLinks returns recently observed links suitable for route calculation.
func (s *LinkStateStore) HealthyLinks(now time.Time, maxAge time.Duration) []Link {
	if maxAge <= 0 {
		return []Link{}
	}
	if now.IsZero() {
		now = time.Now().UTC()
	} else {
		now = now.UTC()
	}
	out := make([]Link, 0)
	for _, l := range s.Snapshot() {
		if l.State == LinkUp && !l.UpdatedAt.IsZero() && !now.Before(l.UpdatedAt) && now.Sub(l.UpdatedAt) <= maxAge {
			out = append(out, l)
		}
	}
	return out
}
