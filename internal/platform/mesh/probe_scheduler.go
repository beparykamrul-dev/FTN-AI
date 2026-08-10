package mesh

import (
	"context"
	"sync"
	"time"
)

type ProbeFunc func(context.Context, Link) ProbeResult

type ProbeScheduler struct {
	store *LinkStateStore
	interval time.Duration
	probe ProbeFunc
}

func NewProbeScheduler(store *LinkStateStore, interval time.Duration, probe ProbeFunc) *ProbeScheduler {
	if interval <= 0 { interval = 5 * time.Second }
	return &ProbeScheduler{store: store, interval: interval, probe: probe}
}

// Run periodically probes a snapshot of links and updates their observed state.
// The scheduler is transport-neutral; ICMP, TCP, agent or synthetic probes can
// be supplied by the caller.
func (s *ProbeScheduler) Run(ctx context.Context) {
	if s.store == nil || s.probe == nil { return }
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		s.probeAll(ctx)
		select {
		case <-ctx.Done(): return
		case <-ticker.C:
		}
	}
}

func (s *ProbeScheduler) probeAll(ctx context.Context) {
	links := s.store.Snapshot()
	var wg sync.WaitGroup
	for _, link := range links {
		link := link
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ctx.Err() != nil { return }
			r := s.probe(ctx, link)
			link.LatencyMS = r.LatencyMS
			link.LossPercent = r.LossPercent
			link.State = EvaluateProbe(r, HealthThresholds{MaxLatencyMS: 100, MaxLossPercent: 5})
			link.UpdatedAt = r.ObservedAt.UTC()
			s.store.Upsert(link)
		}()
	}
	wg.Wait()
}
