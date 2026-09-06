package router

import (
	"context"
	"testing"
	"time"
)

type nextRegistry struct{ services []Service }
func (r nextRegistry) Services(context.Context, string) ([]Service, error) { return r.services, nil }

func TestResolvePrefersHealthyPreferredRegion(t *testing.T) {
	r := New(nextRegistry{services: []Service{
		{ID:"b", Endpoint:"https://b", Region:"other", Healthy:true, RTT:time.Millisecond},
		{ID:"a", Endpoint:"https://a", Region:"preferred", Healthy:true, RTT:50*time.Millisecond},
	}}, time.Minute)
	got, err := r.Resolve(context.Background(), "svc", "preferred")
	if err != nil { t.Fatal(err) }
	if got.ID != "a" { t.Fatalf("route=%q, want preferred region", got.ServiceID) }
}

func TestResolveFailsClosedOnNilContext(t *testing.T) {
	r := New(nextRegistry{}, time.Second)
	if _, err := r.Resolve(nil, "svc", ""); err == nil { t.Fatal("nil context must fail") }
}
