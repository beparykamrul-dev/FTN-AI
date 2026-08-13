package integration

import (
	"sync"
	"time"
)

// DNSMetric is a lightweight observation exported to FTN Metrics/SDN.
type DNSMetric struct {
	Provider  string
	Name      string
	Success   bool
	Secure    bool
	Latency   time.Duration
	Timestamp time.Time
}

// MetricsSink keeps the DNS integration dependency-free. A production sink can
// forward snapshots to FTN Metrics, Prometheus, or another configured backend.
type MetricsSink struct {
	mu    sync.RWMutex
	last  []DNSMetric
	limit int
}

func NewMetricsSink(limit int) *MetricsSink {
	if limit <= 0 { limit = 256 }
	return &MetricsSink{limit: limit}
}

func (s *MetricsSink) Record(m DNSMetric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m.Timestamp.IsZero() { m.Timestamp = time.Now() }
	s.last = append(s.last, m)
	if len(s.last) > s.limit { s.last = s.last[len(s.last)-s.limit:] }
}

func (s *MetricsSink) Snapshot() []DNSMetric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]DNSMetric, len(s.last))
	copy(out, s.last)
	return out
}
