package traffic

import (
	"testing"
	"time"
)

func TestDefaultServicesContainCustomerRealtimeTraffic(t *testing.T) {
	services := DefaultServices()
	want := map[string]bool{"whatsapp": true, "telegram": true, "imo": true, "pubg": true, "freefire": true, "realtime-generic": true}
	for _, s := range services { delete(want, s.ID) }
	if len(want) != 0 { t.Fatalf("missing services: %#v", want) }
}

func TestSelectPrefersHealthyLowJitterPath(t *testing.T) {
	now := time.Now()
	service := Service{ID: "pubg", Class: ClassGaming, Priority: 95, DSCP: 46}
	got, ok := Select([]Observation{
		{PathID: "slow", ServiceID: "pubg", Class: ClassGaming, LatencyMs: 50, JitterMs: 12, PacketLoss: 1, Congestion: .2, Healthy: true, ObservedAt: now},
		{PathID: "fast", ServiceID: "pubg", Class: ClassGaming, LatencyMs: 20, JitterMs: 3, PacketLoss: .1, Congestion: .05, Healthy: true, ObservedAt: now},
	}, service)
	if !ok || got.PathID != "fast" { t.Fatalf("unexpected decision: %#v ok=%v", got, ok) }
	if got.DSCP != 46 || got.HoldDownSec <= 0 { t.Fatalf("missing QoS policy: %#v", got) }
}

func TestSelectRejectsUnhealthyPaths(t *testing.T) {
	service := Service{ID: "telegram", Class: ClassRealtime}
	if _, ok := Select([]Observation{{PathID: "p1", ServiceID: "telegram", Class: ClassRealtime, Healthy: false}}, service); ok { t.Fatal("unhealthy path selected") }
}
