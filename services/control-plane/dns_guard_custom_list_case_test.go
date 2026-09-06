package main

import (
 "encoding/json"
 "testing"
)

func TestDNSGuardCustomListsAreCaseInsensitive(t *testing.T) {
 p:=DNSGuardProfile{CustomBlocklist:json.RawMessage(`["Example.NET"]`)}
 got:=CompileDNSGuardDecision(p,"other","example.net.")
 if got.Decision!="block"{t.Fatalf("unexpected decision: %#v",got)}
}
