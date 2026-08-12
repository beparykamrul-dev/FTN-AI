package dnscontrol

import "testing"

func TestCompareRevision(t *testing.T) {
	if got := CompareRevision("r1", NodeState{NodeID: "n1", ZoneRevision: "r1", Healthy: true}); got != Consistent { t.Fatalf("got %q", got) }
	if got := CompareRevision("r2", NodeState{NodeID: "n1", ZoneRevision: "r1", Healthy: true}); got != Drifted { t.Fatalf("got %q", got) }
	if got := CompareRevision("r1", NodeState{NodeID: "n1", ZoneRevision: "r1", Healthy: false}); got != Unavailable { t.Fatalf("got %q", got) }
}
