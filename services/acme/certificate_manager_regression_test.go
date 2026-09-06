package acme

import (
	"testing"
	"time"
)

func TestCertificateValidationRejectsMissingFields(t *testing.T) {
	if Validate(Certificate{}) == nil {
		t.Fatal("expected invalid certificate to be rejected")
	}
}

func TestCertificateExpiringSoonRejectsInvalidWindow(t *testing.T) {
	c := Certificate{ID: "id", Subject: "subject", Issuer: "issuer", ExpiresAt: time.Now().Add(time.Hour)}
	if IsExpiringSoon(c, time.Now(), -time.Second) {
		t.Fatal("negative window must not report expiring")
	}
}
