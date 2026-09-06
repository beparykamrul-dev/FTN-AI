package ftnmesh

import "testing"

func TestMeshSnapshotRejectsUnknownLinkEndpoint(t *testing.T) {
	s := MeshSnapshot{Nodes: []Node{{ID: "n1"}}, Links: []Link{{ID: "l1", A: "n1", B: "missing", CapacityMbps: 1}}}
	if s.Validate() == nil { t.Fatal("link endpoint outside node set must be rejected") }
}

func TestMeshSnapshotRejectsDuplicateRouteTarget(t *testing.T) {
	s := MeshSnapshot{Nodes: []Node{{ID: "n1"}, {ID: "n2"}}, Routes: []RouteIntent{{ID: "r1", Prefix: "10.0.0.0/24", Source: "n1", Targets: []string{"n2", "n2"}}}}
	if s.Validate() == nil { t.Fatal("duplicate route target must be rejected") }
}
