package proxy

import "testing"

func TestDefaultSecurityPolicyIsValid(t *testing.T) {
	p := DefaultSecurityPolicy()
	if !p.Valid() { t.Fatal("default security policy must be valid") }
	if !p.RequireTLS || !p.RejectPrivateTargets { t.Fatal("secure defaults must require TLS and reject private targets") }
}
