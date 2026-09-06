package builder

import "testing"

func TestManifestRejectsCRLFEnvironment(t *testing.T) {
	m := Manifest{Name: "app", Template: "unknown", Environment: map[string]string{"KEY\n": "value"}}
	if m.Validate() == nil { t.Fatal("environment key with newline must be rejected") }
}

func TestManifestNormalizationDeduplicatesFeatures(t *testing.T) {
	m := Manifest{Name: "app", Template: "unknown", Features: []string{"z", "a", "a"}}
	n := m.Normalize()
	if len(n.Features) != 2 || n.Features[0] != "a" || n.Features[1] != "z" { t.Fatalf("unexpected normalized features: %#v", n.Features) }
}
