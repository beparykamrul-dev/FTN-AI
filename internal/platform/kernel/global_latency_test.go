package kernel

import (
	"testing"
	"time"
)

func TestGlobalPathRouterPrefersHealthyLowScorePath(t *testing.T) {
	r := NewGlobalPathRouter(time.Minute)
	now := time.Now()
	r.Observe(GlobalPathSample{Source: "pop-a", Destination: "pop-b", RTT: 20 * time.Millisecond, Jitter: 2 * time.Millisecond, Loss: 0.01, Utilization: 0.20, Healthy: true, ObservedAt: now})
	r.Observe(GlobalPathSample{Source: "pop-a", Destination: "pop-b", RTT: 60 * time.Millisecond, Jitter: 1 * time.Millisecond, Loss: 0, Utilization: 0.10, Healthy: true, ObservedAt: now})

	c := r.Candidates("pop-a", "pop-b", now)
	if len(c) != 2 || c[0].Sample.RTT != 20*time.Millisecond {
		t.Fatalf("unexpected path ordering: %#v", c)
	}
}

func TestGlobalPathRouterDropsStaleAndUnhealthyPaths(t *testing.T) {
	r := NewGlobalPathRouter(time.Second)
	now := time.Now()
	r.Observe(GlobalPathSample{Source: "a", Destination: "b", RTT: 10 * time.Millisecond, Healthy: false, ObservedAt: now})
	r.Observe(GlobalPathSample{Source: "a", Destination: "b", RTT: 11 * time.Millisecond, Healthy: true, ObservedAt: now.Add(-2 * time.Second)})
	if got := r.Candidates("a", "b", now); len(got) != 0 {
		t.Fatalf("expected no usable paths, got %d", len(got))
	}
}
