package ftnmesh

import "testing"

func TestMeshSnapshotRejectsSelfLink(t *testing.T) {
	s := MeshSnapshot{
		Nodes: []Node{{ID: "n1"}},
		Links: []Link{{ID: "l1", A: "n1", B: "n1", CapacityMbps: 100}},
	}
	if err := s.Validate(); err == nil {
		t.Fatal("self-link must be rejected")
	}
}

func TestMeshSnapshotRejectsUnknownEndpoint(t *testing.T) {
	s := MeshSnapshot{
		Nodes: []Node{{ID: "n1"}},
		Links: []Link{{ID: "l1", A: "n1", B: "n2", CapacityMbps: 100}},
	}
	if err := s.Validate(); err == nil {
		t.Fatal("unknown endpoint must be rejected")
	}
}
