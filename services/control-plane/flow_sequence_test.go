package main

import "testing"

func TestFlowSequenceTracker(t *testing.T) {
	tk := NewFlowSequenceTracker()
	if g := tk.Observe("192.0.2.1", 9, 100, 1); g != 0 { t.Fatalf("initial gap=%d", g) }
	if g := tk.Observe("192.0.2.1", 9, 101, 1); g != 0 { t.Fatalf("contiguous gap=%d", g) }
	if g := tk.Observe("192.0.2.1", 9, 105, 1); g != 3 { t.Fatalf("gap=%d want 3", g) }
	if got := tk.State("192.0.2.1", 9); got.Gaps != 3 || got.Packets != 3 { t.Fatalf("state=%+v", got) }
}

func TestFlowSequenceTrackerExporterIsolation(t *testing.T) {
	tk := NewFlowSequenceTracker()
	tk.Observe("192.0.2.1", 10, 10, 1)
	if got := tk.State("192.0.2.2", 10); got.Initialized { t.Fatal("state leaked across exporters") }
}
