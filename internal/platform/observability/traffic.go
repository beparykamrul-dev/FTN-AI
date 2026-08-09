package observability

import (
	"sync"
	"time"
)

type TrafficSample struct {
	Interface string    `json:"interface"`
	SourceIP  string    `json:"source_ip,omitempty"`
	DestIP    string    `json:"dest_ip,omitempty"`
	Protocol  string    `json:"protocol,omitempty"`
	App       string    `json:"application,omitempty"`
	Bytes     uint64    `json:"bytes"`
	Packets   uint64    `json:"packets"`
	At        time.Time `json:"at"`
}

type TrafficStore struct { mu sync.RWMutex; samples []TrafficSample; limit int }

func NewTrafficStore(limit int) *TrafficStore { if limit < 1 { limit = 10000 }; return &TrafficStore{limit: limit} }

func (s *TrafficStore) Add(v TrafficSample) {
	s.mu.Lock(); defer s.mu.Unlock()
	if v.At.IsZero() { v.At = time.Now().UTC() }
	s.samples = append(s.samples, v)
	if len(s.samples) > s.limit { s.samples = s.samples[len(s.samples)-s.limit:] }
}

func (s *TrafficStore) List() []TrafficSample {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]TrafficSample, len(s.samples)); copy(out, s.samples); return out
}
