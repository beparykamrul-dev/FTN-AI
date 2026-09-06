package module

import "testing"

func TestManifestRejectsSelfDependency(t *testing.T) {
	m := Manifest{Name: "core", Version: "1.0.0", Dependencies: []string{"core"}}
	if m.Validate() == nil { t.Fatal("module self-dependency must be rejected") }
}

func TestManifestRejectsDuplicateCapability(t *testing.T) {
	m := Manifest{Name: "core", Version: "1.0.0", Capabilities: []string{"read", "read"}}
	if m.Validate() == nil { t.Fatal("duplicate capability must be rejected") }
}
