package router

import (
	"context"
	"testing"
	"time"
)

type testRegistry struct{ services []Service }

func (t testRegistry) Services(context.Context, string) ([]Service, error) {
	return t.services, nil
}

func TestResolvePrefersHealthyPreferredRegionThenLatency(t *testing.T) {
	r := New(testRegistry{services: []Service{
		{ID: "remote", Endpoint: "https://remote", Region: "remote", Healthy: true, RTT: 5 * time.Millisecond},
		{ID: "local", Endpoint: "https://local", Region: "local", Healthy: true, RTT: 30 * time.Millisecond},
	}}, time.Second)

	got, err := r.Resolve(context.Background(), "dns", "local")
	if err != nil {
		t.Fatal(err)
	}
	if got.ServiceID != "local" {
		t.Fatalf("expected preferred-region service, got %q", got.ServiceID)
	}
}

func TestResolveSkipsUnhealthyRoutes(t *testing.T) {
	r := New(testRegistry{services: []Service{
		{ID: "bad", Endpoint: "https://bad", Healthy: false, RTT: time.Millisecond},
		{ID: "good", Endpoint: "https://good", Healthy: true, RTT: 10 * time.Millisecond},
	}}, time.Second)

	got, err := r.Resolve(context.Background(), "dns", "")
	if err != nil {
		t.Fatal(err)
	}
	if got.ServiceID != "good" {
		t.Fatalf("expected healthy service, got %q", got.ServiceID)
	}
}

func TestResolveReturnsNoRoute(t *testing.T) {
	r := New(testRegistry{}, time.Second)
	if _, err := r.Resolve(context.Background(), "dns", ""); err != ErrNoRoute {
		t.Fatalf("expected ErrNoRoute, got %v", err)
	}
}

func TestResolveUsesCache(t *testing.T) {
	reg := &countingRegistry{services: []Service{{ID: "one", Endpoint: "https://one", Healthy: true}}}
	r := New(reg, time.Minute)
	if _, err := r.Resolve(context.Background(), "dns", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Resolve(context.Background(), "dns", ""); err != nil {
		t.Fatal(err)
	}
	if reg.calls != 1 {
		t.Fatalf("expected one registry lookup, got %d", reg.calls)
	}
}

func TestResolveCacheIsScopedByPreferredRegion(t *testing.T) {
	reg := &countingRegistry{services: []Service{
		{ID: "local", Endpoint: "https://local", Region: "local", Healthy: true, RTT: 30 * time.Millisecond},
		{ID: "remote", Endpoint: "https://remote", Region: "remote", Healthy: true, RTT: 5 * time.Millisecond},
	}}
	r := New(reg, time.Minute)

	local, err := r.Resolve(context.Background(), "dns", "local")
	if err != nil {
		t.Fatal(err)
	}
	remote, err := r.Resolve(context.Background(), "dns", "remote")
	if err != nil {
		t.Fatal(err)
	}
	if local.ServiceID != "local" || remote.ServiceID != "remote" {
		t.Fatalf("cache crossed region boundary: local=%q remote=%q", local.ServiceID, remote.ServiceID)
	}
	if reg.calls != 2 {
		t.Fatalf("expected separate registry lookups per region, got %d", reg.calls)
	}
}

type countingRegistry struct {
	services []Service
	calls    int
}

func (c *countingRegistry) Services(context.Context, string) ([]Service, error) {
	c.calls++
	return c.services, nil
}
