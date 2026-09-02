package ftnmesh

import "testing"

func TestMeshSnapshotValidate(t *testing.T) {
	valid := MeshSnapshot{
		Nodes: []Node{{ID: "core-a"}, {ID: "pop-a"}},
		Links: []Link{{ID: "l1", A: "core-a", B: "pop-a", CapacityMbps: 10000}},
		Routes: []RouteIntent{{ID: "r1", Prefix: "203.0.113.0/24", Source: "core-a", Targets: []string{"pop-a"}, Approved: true}},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid snapshot rejected: %v", err)
	}
}

func TestMeshSnapshotRejectsDuplicateNodes(t *testing.T) {
	s := MeshSnapshot{Nodes: []Node{{ID: "core-a"}, {ID: "core-a"}}}
	if err := s.Validate(); err == nil {
		t.Fatal("duplicate node IDs were accepted")
	}
}

func TestMeshSnapshotRejectsSelfLinkAndZeroCapacity(t *testing.T) {
	cases := []MeshSnapshot{
		{Links: []Link{{ID: "l1", A: "a", B: "a", CapacityMbps: 1000}}},
		{Links: []Link{{ID: "l1", A: "a", B: "b"}}},
	}
	for i, s := range cases {
		if err := s.Validate(); err == nil {
			t.Fatalf("case %d was accepted", i)
		}
	}
}

func TestMeshSnapshotRequiresTargetForApprovedRoute(t *testing.T) {
	s := MeshSnapshot{Routes: []RouteIntent{{ID: "r1", Prefix: "203.0.113.0/24", Source: "core-a", Approved: true}}}
	if err := s.Validate(); err == nil {
		t.Fatal("approved route without target was accepted")
	}
}
