package main

import "testing"

func TestDNSGuardEmptyDomainDoesNotMatch(t *testing.T) {
 got:=CompileDNSGuardDecision(DNSGuardProfile{BlockAds:true},"ads","")
 if got.Decision!="block" { t.Fatalf("category policy should still apply: %#v",got) }
}
