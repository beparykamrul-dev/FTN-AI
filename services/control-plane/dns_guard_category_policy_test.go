package main

import "testing"

func TestDNSGuardCategoryPolicy(t *testing.T) {
 p:=DNSGuardProfile{BlockMalware:true}
 got:=CompileDNSGuardDecision(p," MALWARE ","safe.example")
 if got.Decision!="block" || got.Reason!="category_policy" { t.Fatalf("unexpected decision: %#v",got) }
}
