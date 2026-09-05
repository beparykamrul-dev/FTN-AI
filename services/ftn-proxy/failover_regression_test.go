package proxy

import "testing"

func TestDefaultFailoverPolicyIsValid(t *testing.T) {
	p := DefaultFailoverPolicy()
	if !p.Valid() { t.Fatal("default failover policy must be valid") }
	if p.MaxAttempts <= 0 || p.RetryDelay <= 0 { t.Fatal("failover retry bounds must be positive") }
}
