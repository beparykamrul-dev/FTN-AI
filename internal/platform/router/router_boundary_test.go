package router

import (
	"context"
	"testing"
	"time"
)

type staticRegistry struct{ items []Service }

func (s staticRegistry) Services(context.Context, string) ([]Service, error) { return s.items, nil }

func TestRouterPrefersHealthyPreferredRegion(t *testing.T) {
	r := New(staticRegistry{items: []Service{
		{ID: "b", Endpoint: "b", Region: "chittagong", Healthy: true, RTT: 5 * time.Millisecond},
		{ID: "a", Endpoint: "a", Region: "dhaka", Healthy: true, RTT: 20 * time.Millisecond},
	}}, time.Second)
	got, err := r.Resolve(context.Background(), "internet", "dhaka")
	if err != nil || got.ID != "a" {
		t.Fatalf("unexpected route: %+v err=%v", got, err)
	}
}
