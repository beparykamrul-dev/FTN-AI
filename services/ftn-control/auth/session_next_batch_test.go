package auth

import "testing"

func TestNewSessionRejectsEmptyIdentity(t *testing.T) {
	if _, err := NewSession(" "); err == nil { t.Fatal("empty identity must fail") }
	if s, err := NewSession("identity-1"); err != nil || s.ID == "" || s.IdentityID != "identity-1" || s.Revoked { t.Fatalf("unexpected session: %#v err=%v", s, err) }
}
