package proxy

import "testing"

func TestDefaultUserSecurityPolicyRejectsUnsafeRequest(t *testing.T) {
	p := DefaultUserSecurityPolicy()
	if p.ValidateUserRequest("http", true, true, false, true) { t.Fatal("HTTP must be rejected") }
	if p.ValidateUserRequest("https", true, true, true, true) { t.Fatal("replayed request must be rejected") }
	if !p.ValidateUserRequest("https", true, true, false, true) { t.Fatal("secure request should pass") }
}
