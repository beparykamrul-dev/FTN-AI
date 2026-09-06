package main

import (
 "encoding/json"
 "testing"
)

func TestDNSGuardInvalidAllowlistIsIgnored(t *testing.T) {
 p:=DNSGuardProfile{CustomAllowlist:json.RawMessage(`{"bad":true}`),BlockAds:true}
 got:=CompileDNSGuardDecision(p,"ads","example.com")
 if got.Decision!="block" { t.Fatalf("invalid allowlist must not disable policy: %#v",got) }
}
