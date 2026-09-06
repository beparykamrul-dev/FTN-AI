package proxy

import (
	"testing"
	"time"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	p := DefaultPolicy()
	if !p.Valid() { t.Fatal("default proxy policy must be valid") }
	if p.MaxBodyBytes <= 0 { t.Fatal("proxy body limit must be positive") }
}

func TestPolicyRejectsInvalidBounds(t *testing.T) {
	p := DefaultPolicy()
	p.IdleTimeout = time.Nanosecond
	if p.Valid() { t.Fatal("idle timeout below connect timeout must be rejected") }
}

func TestPolicyRejectsPlainHTTPWhenTLSRequired(t *testing.T) {
	if DefaultPolicy().Validate("http", 1) { t.Fatal("plain HTTP must be rejected when TLS is required") }
}
