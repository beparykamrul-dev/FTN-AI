package proxy

import "testing"

func TestPolicyValidateSchemeAndBody(t *testing.T) {
	p := DefaultPolicy()
	if !p.Validate("https", 1024) { t.Fatal("expected valid HTTPS request") }
	if p.Validate("http", 1024) { t.Fatal("expected TLS policy to reject HTTP") }
	if p.Validate("ftp", 1024) { t.Fatal("expected unsupported scheme to be rejected") }
	if p.Validate("https", -1) { t.Fatal("expected negative body size to be rejected") }
	if p.Validate("https", p.MaxBodyBytes+1) { t.Fatal("expected oversized body to be rejected") }
}
