package main

import "testing"

func TestDNSGuardReasonsAreStable(t *testing.T) {
 cases:=[]struct{p DNSGuardProfile;cat,domain,reason string}{
  {DNSGuardProfile{BlockAdult:true},"adult","x","category_policy"},
  {DNSGuardProfile{},"news","x","default_allow"},
 }
 for _,tc:=range cases{got:=CompileDNSGuardDecision(tc.p,tc.cat,tc.domain);if got.Reason!=tc.reason{t.Fatalf("got %q want %q",got.Reason,tc.reason)}}
}
