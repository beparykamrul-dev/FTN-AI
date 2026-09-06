package main

import "testing"

func TestDNSGuardDefaultAllow(t *testing.T) {
 got:=CompileDNSGuardDecision(DNSGuardProfile{},"unknown","example.com")
 if got.Decision!="allow" || got.Reason!="default_allow" { t.Fatalf("unexpected decision: %#v",got) }
}
