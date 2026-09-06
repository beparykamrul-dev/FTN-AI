package inventory

import "testing"

func TestNetBoxClientAndObjectBoundaries(t *testing.T) {
 if err := (NetBoxClient{BaseURL: "https://netbox.example", TokenRef: "secret-ref"}).Validate(); err != nil { t.Fatalf("valid NetBox client rejected: %v", err) }
 if err := (NetBoxClient{BaseURL: "http://netbox.example", TokenRef: "secret-ref"}).Validate(); err == nil { t.Fatal("HTTP NetBox endpoint accepted") }
 o := (NetBoxObject{ID: " 42 ", Type: " DEVICE ", Name: " core "}).Normalize(); if !o.Valid() || o.ID != "42" || o.Type != "device" || o.Name != "core" { t.Fatalf("normalized object=%+v", o) }
}
