package module

import "testing"

func TestManifestValidationAndNormalization(t *testing.T) {
 m := Manifest{Name: " mod ", Version: " 1.0 ", Capabilities: []string{"dns", "dns", " network "}, Dependencies: []string{"core", "core"}}
 if err := m.Validate(); err != nil { t.Fatalf("valid manifest rejected: %v", err) }
 n := m.Normalize(); if n.Name != "mod" || n.Version != "1.0" || len(n.Capabilities) != 2 || len(n.Dependencies) != 1 { t.Fatalf("normalized manifest=%+v", n) }
 m.Dependencies = []string{"mod"}; if err := m.Validate(); err == nil { t.Fatal("self dependency accepted") }
}
