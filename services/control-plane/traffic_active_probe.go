package main

import (
	"context"
	"errors"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	trafficProbeDefaultInterval = 15 * time.Second
	trafficProbeDefaultTimeout  = 2 * time.Second
	trafficProbeWindow          = 8
	trafficProbeFailureLimit    = 3
	trafficProbeMaxConcurrency  = 16
)

type TrafficProbeTarget struct {
	ServiceID string
	PathID    string
	Address   string
}

type trafficProbeSample struct {
	Latency time.Duration
	Success bool
	At      time.Time
}

type trafficProbeState struct {
	Samples            []trafficProbeSample
	ConsecutiveFailure int
}

type TrafficActiveProbe struct {
	runtime   *TrafficRuntime
	interval  time.Duration
	timeout   time.Duration
	mu        sync.Mutex
	states    map[string]trafficProbeState
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	started   bool
}

func NewTrafficActiveProbe(runtime *TrafficRuntime, interval, timeout time.Duration) *TrafficActiveProbe {
	if interval <= 0 {
		interval = trafficProbeDefaultInterval
	}
	if timeout <= 0 {
		timeout = trafficProbeDefaultTimeout
	}
	return &TrafficActiveProbe{runtime: runtime, interval: interval, timeout: timeout, states: make(map[string]trafficProbeState)}
}

func (p *TrafficActiveProbe) Start(ctx context.Context, targets func() []TrafficProbeTarget) error {
	if p == nil || p.runtime == nil || targets == nil {
		return errors.New("traffic_probe_configuration_required")
	}
	p.mu.Lock()
	if p.started {
		p.mu.Unlock()
		return errors.New("traffic_probe_already_started")
	}
	probeCtx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.started = true
	p.mu.Unlock()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.run(probeCtx, targets)
	}()
	return nil
}

func (p *TrafficActiveProbe) Close() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	cancel := p.cancel
	p.cancel = nil
	p.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	p.wg.Wait()
	return nil
}

func (p *TrafficActiveProbe) run(ctx context.Context, targets func() []TrafficProbeTarget) {
	p.probeOnce(ctx, targets())
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.probeOnce(ctx, targets())
		}
	}
}

func (p *TrafficActiveProbe) probeOnce(ctx context.Context, targets []TrafficProbeTarget) {
	clean := make([]TrafficProbeTarget, 0, len(targets))
	seen := make(map[string]struct{}, len(targets))
	for _, target := range targets {
		target.ServiceID = strings.TrimSpace(target.ServiceID)
		target.PathID = strings.TrimSpace(target.PathID)
		target.Address = strings.TrimSpace(target.Address)
		if target.ServiceID == "" || target.PathID == "" || target.Address == "" {
			continue
		}
		if _, ok := trafficPolicyByID(target.ServiceID); !ok {
			continue
		}
		key := trafficQualityKey(target.ServiceID, target.PathID)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		clean = append(clean, target)
	}
	sort.Slice(clean, func(i, j int) bool {
		if clean[i].ServiceID == clean[j].ServiceID {
			return clean[i].PathID < clean[j].PathID
		}
		return clean[i].ServiceID < clean[j].ServiceID
	})
	if len(clean) == 0 {
		return
	}

	sem := make(chan struct{}, trafficProbeMaxConcurrency)
	var wg sync.WaitGroup
	for _, target := range clean {
		target := target
		select {
		case <-ctx.Done():
			return
		default:
		}
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			p.probeTarget(ctx, target)
		}()
	}
	wg.Wait()
}

func (p *TrafficActiveProbe) probeTarget(ctx context.Context, target TrafficProbeTarget) {
	started := time.Now()
	dialer := net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", target.Address)
	latency := time.Since(started)
	if conn != nil {
		_ = conn.Close()
	}
	now := time.Now().UTC()

	key := trafficQualityKey(target.ServiceID, target.PathID)
	p.mu.Lock()
	state := p.states[key]
	sample := trafficProbeSample{Latency: latency, Success: err == nil, At: now}
	state.Samples = append(state.Samples, sample)
	if len(state.Samples) > trafficProbeWindow {
		state.Samples = state.Samples[len(state.Samples)-trafficProbeWindow:]
	}
	if sample.Success {
		state.ConsecutiveFailure = 0
	} else {
		state.ConsecutiveFailure++
	}
	p.states[key] = state
	observation := summarizeTrafficProbe(target, state, now)
	p.mu.Unlock()

	_ = p.runtime.UpsertQuality(observation, now)
}

func summarizeTrafficProbe(target TrafficProbeTarget, state trafficProbeState, now time.Time) TrafficQualityObservation {
	var success int
	var totalLatency time.Duration
	var previous time.Duration
	var jitter time.Duration
	var previousSet bool
	for _, sample := range state.Samples {
		if !sample.Success {
			continue
		}
		success++
		totalLatency += sample.Latency
		if previousSet {
			delta := sample.Latency - previous
			if delta < 0 {
				delta = -delta
			}
			jitter += delta
		}
		previous = sample.Latency
		previousSet = true
	}
	loss := 100.0
	if len(state.Samples) > 0 {
		loss = float64(len(state.Samples)-success) * 100 / float64(len(state.Samples))
	}
	latencyMs := 0.0
	jitterMs := 0.0
	if success > 0 {
		latencyMs = totalLatency.Seconds() * 1000 / float64(success)
		if success > 1 {
			jitterMs = jitter.Seconds() * 1000 / float64(success-1)
		}
	}
	healthy := success > 0 && state.ConsecutiveFailure < trafficProbeFailureLimit
	return TrafficQualityObservation{
		PathID:      target.PathID,
		ServiceID:   target.ServiceID,
		Class:       trafficPolicyClass(target.ServiceID),
		LatencyMs:   latencyMs,
		JitterMs:    jitterMs,
		PacketLoss:  loss,
		Congestion:  0,
		Healthy:     healthy,
		ObservedAt:  now,
	}
}

func trafficPolicyClass(serviceID string) TrafficClass {
	if p, ok := trafficPolicyByID(serviceID); ok {
		return p.Class
	}
	return TrafficNormal
}
