package dns

import "testing"

func TestNodeRegistryRejectsWrongZone(t *testing.T) {
	r := NewNodeRegistry()
	if r.Upsert(DNSNode{ID: "n1", Name: "node", ProviderID: "p1", Scope: NodeScopeGlobal, Zone: "example.com", Enabled: true}) {
		t.Fatal("expected wrong zone to be rejected")
	}
}

func TestNodeRegistryAcceptsAdditiveFamilyTimeNetNode(t *testing.T) {
	r := NewNodeRegistry()
	if !r.Upsert(DNSNode{ID: "n1", Name: "node", ProviderID: "p1", Scope: NodeScopeGlobal, Zone: FamilyTimeNetZone, Enabled: true}) {
		t.Fatal("expected valid node")
	}
	if _, ok := r.Get("n1"); !ok {
		t.Fatal("expected node in registry")
	}
}

func TestNodeRegistryListsByScope(t *testing.T) {
	r := NewNodeRegistry()
	_ = r.Upsert(DNSNode{ID: "g", Name: "global", ProviderID: "p", Scope: NodeScopeGlobal, Zone: FamilyTimeNetZone, Enabled: true})
	_ = r.Upsert(DNSNode{ID: "l", Name: "local", ProviderID: "p", Scope: NodeScopeLocal, Zone: FamilyTimeNetZone, Enabled: true})
	if got := r.List(NodeScopeGlobal); len(got) != 1 || got[0].ID != "g" {
		t.Fatalf("unexpected global nodes: %+v", got)
	}
	if got := r.List(NodeScopeLocal); len(got) != 1 || got[0].ID != "l" {
		t.Fatalf("unexpected local nodes: %+v", got)
	}
}
