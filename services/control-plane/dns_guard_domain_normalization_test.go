package main

import (
 "encoding/json"
 "testing"
)

func TestDNSGuardDomainNormalization(t *testing.T) {
 p:=DNSGuardProfile{CustomAllowlist:json.RawMessage(`["Example.COM."]`)}
 got:=CompileDNSGuardDecision(p,"ads"," example.com. ")
 if got.Decision!="allow" || got.Reason!="custom_allowlist" { t.Fatalf("unexpected decision: %#v",got) }
}
