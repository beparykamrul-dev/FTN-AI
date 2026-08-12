package authenticator

import (
	"context"
	"testing"
	"time"
)

type testVerifier struct{}

func (testVerifier) Verify(context.Context, string, string) error { return nil }

type testAudit struct{ events []AuditEvent }

func (a *testAudit) Record(_ context.Context, event AuditEvent) { a.events = append(a.events, event) }

func TestAuthenticatorSessionRoundTrip(t *testing.T) {
	audit := &testAudit{}
	a, err := New(testVerifier{}, audit, []byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	a.now = func() time.Time { return now }
	token, err := a.SignSession(Session{ID: "sess-1", IdentityID: "id-1", IssuedAt: now, ExpiresAt: now.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if err := a.VerifySession(token); err != nil {
		t.Fatal(err)
	}
	if err := a.VerifySession(token + "x"); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}

func TestAuthenticatorRejectsExpiredSession(t *testing.T) {
	a, err := New(testVerifier{}, nil, []byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	a.now = func() time.Time { return now }
	token, err := a.SignSession(Session{ID: "sess-2", IdentityID: "id-2", IssuedAt: now.Add(-2 * time.Hour), ExpiresAt: now.Add(-time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if err := a.VerifySession(token); err != ErrSessionExpired {
		t.Fatalf("expected ErrSessionExpired, got %v", err)
	}
}

func TestAuthenticatorTOTP(t *testing.T) {
	a, err := New(testVerifier{}, nil, []byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}
	secret := []byte("12345678901234567890")
	code := totp(secret, 0)
	if err := a.VerifyTOTP(secret, code, 0); err != nil {
		t.Fatal(err)
	}
	if err := a.VerifyTOTP(secret, "000000", 0); err == nil {
		t.Fatal("expected invalid TOTP to be rejected")
	}
}

func TestAuthenticatorCredentialAudit(t *testing.T) {
	audit := &testAudit{}
	a, err := New(testVerifier{}, audit, []byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}
	if err := a.VerifyCredential(context.Background(), "id-3", "credential", "req-1"); err != nil {
		t.Fatal(err)
	}
	if len(audit.events) != 1 || audit.events[0].Result != "allowed" {
		t.Fatalf("unexpected audit events: %#v", audit.events)
	}
}
