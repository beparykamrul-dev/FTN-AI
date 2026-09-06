package main

import "testing"

func TestDNSGuardCategoryMatrix(t *testing.T) {
 cases:=[]struct{name string;p DNSGuardProfile;cat string;want string}{
  {"ads",DNSGuardProfile{BlockAds:true},"ads","block"},
  {"trackers",DNSGuardProfile{BlockTrackers:true},"trackers","block"},
  {"phishing",DNSGuardProfile{BlockPhishing:true},"phishing","block"},
  {"gambling",DNSGuardProfile{BlockGambling:true},"gambling","block"},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){if got:=CompileDNSGuardDecision(tc.p,tc.cat,"x.example").Decision;got!=tc.want{t.Fatalf("got %s want %s",got,tc.want)}})}
}
