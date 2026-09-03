package main

import "testing"

func TestRouteEventCorrelationIDDeterministic(t *testing.T) {
	e := FTNBGPEvent{
		Type: "bgp.peer.state",
		Peer: "192.0.2.2",
		State: "established",
		Timestamp: "2026-09-03T00:00:00Z",
	}
	first := routeEventCorrelationID(e)
	second := routeEventCorrelationID(e)
	if first == "" || first != second {
		t.Fatalf("correlation id is not deterministic: %q %q", first, second)
	}
	if len(first) != len("ftn-routed-")+64 {
		t.Fatalf("unexpected correlation id length: %d", len(first))
	}
}

func TestRouteEventCorrelationIDSeparatesEvents(t *testing.T) {
	base := FTNBGPEvent{
		Type: "bgp.peer.state",
		Peer: "192.0.2.2",
		State: "established",
		Timestamp: "2026-09-03T00:00:00Z",
	}
	changed := base
	changed.State = "idle"
	if routeEventCorrelationID(base) == routeEventCorrelationID(changed) {
		t.Fatal("different routing states must not share a correlation id")
	}
}
