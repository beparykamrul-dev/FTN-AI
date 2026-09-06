package deploy

import "testing"

func TestTargetValidationIdentityBoundaries(t *testing.T) {
 t1 := Target{ID: " t1 ", Name: " core ", IP: "192.0.2.1", Serial: " s1 "}
 if err := t1.Validate(); err != nil { t.Fatalf("valid target rejected: %v", err) }
 n := t1.Normalize(); if n.ID != "t1" || n.Name != "core" || n.Serial != "s1" { t.Fatalf("normalized target=%+v", n) }
 if err := (Target{ID: "t1", Name: "core", IP: "bad", Serial: "s1"}).Validate(); err == nil { t.Fatal("invalid target IP accepted") }
 if err := (Target{ID: "t1", Name: "core"}).Validate(); err == nil { t.Fatal("target without serial or agent accepted") }
}
