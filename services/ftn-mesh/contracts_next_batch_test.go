package ftnmesh

import "testing"

func TestMeshSnapshotValidationBoundaries(t *testing.T) {
 s := MeshSnapshot{Nodes: []Node{{ID: "a"}, {ID: "b"}}, Links: []Link{{ID: "l1", A: "a", B: "b", CapacityMbps: 100, LatencyMs: 1, LossPct: 0}}}
 if err := s.Validate(); err != nil { t.Fatalf("valid mesh rejected: %v", err) }
 s.Links[0].LossPct = 101; if err := s.Validate(); err == nil { t.Fatal("invalid link loss accepted") }
 s = MeshSnapshot{Nodes: []Node{{ID: "a"}, {ID: "b"}}, Links: []Link{{ID: "l1", A: "a", B: "c", CapacityMbps: 100}}}; if err := s.Validate(); err == nil { t.Fatal("unknown link endpoint accepted") }
}
