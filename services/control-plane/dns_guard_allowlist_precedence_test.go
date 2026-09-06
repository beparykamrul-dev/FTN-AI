package main

import (
 "encoding/json"
 "testing"
)

func TestDNSGuardAllowlistPrecedence(t *testing.T) {
 p:=DNSGuardProfile{BlockAds:true,CustomAllowlist:json.RawMessage(`["ads.example"]`)}
 got:=CompileDNSGuardDecision(p,"ads","ads.example")
 if got.Decision!="allow" { t.Fatalf("allowlist must win: %#v",got) }
}
