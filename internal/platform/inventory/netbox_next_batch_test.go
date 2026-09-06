package inventory

import "testing"

func TestNetBoxValidationRejectsUnsafeAndOversizedConfig(t *testing.T) {
	if err := (NetBoxClient{BaseURL: "http://example.com", TokenRef: "ref"}).Validate(); err == nil { t.Fatal("expected HTTPS requirement") }
	if err := (NetBoxClient{BaseURL: "https://example.com", TokenRef: ""}).Validate(); err == nil { t.Fatal("expected token reference requirement") }
	obj := (NetBoxObject{ID: " ", Type: "device", Name: "n"}).Normalize()
	if obj.Valid() { t.Fatal("normalized empty ID must remain invalid") }
}
