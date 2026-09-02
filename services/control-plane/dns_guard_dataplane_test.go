package main

import "testing"

func TestNormalizeDNSName(t *testing.T) {
    if got := NormalizeDNSName(" Example.COM. "); got != "example.com" { t.Fatalf("got %q", got) }
    if got := NormalizeDNSName("127.0.0.1"); got != "" { t.Fatalf("IP must not be treated as domain: %q", got) }
    if got := NormalizeDNSName("bad_label.example"); got != "" { t.Fatalf("invalid label accepted: %q", got) }
}

func TestDNSGuardDataplane(t *testing.T) {
    p := DNSGuardProfile{BlockMalware:true, BlockAds:true}
    got := CompileDNSGuardDataplane(p, DNSGuardDataplaneRequest{Domain:"bad.example.", Category:"malware"})
    if got.Decision != "block" || got.Reason != "category_policy" || got.DomainHash == "" { t.Fatalf("unexpected result: %+v", got) }
}
