package main

import (
	"encoding/json"
	"testing"
)

func TestCompileDNSGuardDecision(t *testing.T) {
	p := DNSGuardProfile{BlockMalware:true, BlockAds:true, CustomAllowlist: json.RawMessage(`["safe.example"]`), CustomBlocklist: json.RawMessage(`["blocked.example"]`)}
	cases := []struct{name,domain,category,decision,reason string}{
		{"explicit allow","safe.example","malware","allow","custom_allowlist"},
		{"explicit block","blocked.example","ads","block","custom_blocklist"},
		{"category block","other.example","malware","block","category_policy"},
		{"default allow","other.example","news","allow","default_allow"},
	}
	for _, tc := range cases { t.Run(tc.name, func(t *testing.T) { got:=CompileDNSGuardDecision(p,tc.category,tc.domain); if got.Decision!=tc.decision||got.Reason!=tc.reason { t.Fatalf("got %+v, want decision=%s reason=%s",got,tc.decision,tc.reason) } }) }
}
