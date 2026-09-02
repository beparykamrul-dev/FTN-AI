package upstream

import (
	"testing"
	"time"
)

func TestMetricsSnapshotIsolation(t *testing.T) {
	m := NewMetrics()
	now := time.Now().UTC()
	m.Set("primary", Snapshot{Healthy: true, Latency: 12 * time.Millisecond, PrefixCount: 80, PrefixLimit: 100, LastChecked: now, State: "up"})

	s, ok := m.Get("primary")
	if !ok || !s.Healthy || s.PrefixCount != 80 || s.PrefixLimit != 100 {
		t.Fatalf("unexpected snapshot: %#v", s)
	}

	all := m.All()
	all["primary"] = Snapshot{}
	s, _ = m.Get("primary")
	if !s.Healthy { t.Fatal("All returned internal mutable state") }
}
