package kernel

import (
	"testing"
	"time"
)

func TestLatencyRouterChoosesLowestHealthyRTT(t *testing.T) {
	r := NewLatencyRouter()
	r.Observe(EndpointSample{Endpoint: "edge-a", RTT: 40 * time.Millisecond, Loss: 0.01, Healthy: true})
	r.Observe(EndpointSample{Endpoint: "edge-b", RTT: 12 * time.Millisecond, Loss: 0.02, Healthy: true})
	r.Observe(EndpointSample{Endpoint: "edge-c", RTT: 1 * time.Millisecond, Loss: 1, Healthy: false})

	best, ok := r.Best()
	if !ok || best.Endpoint != "edge-b" {
		t.Fatalf("expected edge-b, got %+v", best)
	}
}

func TestLatencyRouterIgnoresInvalidSamples(t *testing.T) {
	r := NewLatencyRouter()
	r.Observe(EndpointSample{Endpoint: "bad", RTT: -1, Healthy: true})
	if _, ok := r.Best(); ok {
		t.Fatal("invalid sample must not become a route")
	}
}
