package main

import (
 "encoding/json"
 "testing"
)

func TestDNSGuardBlocklistPrecedence(t *testing.T) {
 p:=DNSGuardProfile{CustomBlocklist:json.RawMessage(`["blocked.example"]`)}
 got:=CompileDNSGuardDecision(p,"other","blocked.example")
 if got.Decision!="block" || got.Reason!="custom_blocklist" { t.Fatalf("unexpected decision: %#v",got) }
}
