package main

import "testing"

func TestDNSGuardMultipleCategoriesRemainIndependent(t *testing.T) {
 p:=DNSGuardProfile{BlockAds:true,BlockMalware:false}
 if got:=CompileDNSGuardDecision(p,"ads","x").Decision;got!="block"{t.Fatalf("ads got %s",got)}
 if got:=CompileDNSGuardDecision(p,"malware","x").Decision;got!="allow"{t.Fatalf("malware got %s",got)}
}
