package proxy

import "testing"

func TestDefaultPolicyIsValid(t *testing.T) {
	p := DefaultPolicy()
	if !p.Valid() { t.Fatal("default proxy policy must be valid") }
	if p.MaxBodyBytes <= 0 { t.Fatal("proxy body limit must be positive") }
}
