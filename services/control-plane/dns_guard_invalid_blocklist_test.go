package main

import (
 "encoding/json"
 "testing"
)

func TestDNSGuardInvalidBlocklistIsIgnored(t *testing.T) {
 p:=DNSGuardProfile{CustomBlocklist:json.RawMessage(`not-json`)}
 got:=CompileDNSGuardDecision(p,"ads","example.com")
 if got.Decision!="allow" { t.Fatalf("malformed blocklist must be ignored: %#v",got) }
}
