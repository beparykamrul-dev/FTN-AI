package acme

import (
 "testing"
 "time"
)

func TestCertificateValidationAndExpiryWindow(t *testing.T) {
 now := time.Now()
 c := Certificate{ID: "c1", Subject: "example.com", Issuer: "FTN-CA", ExpiresAt: now.Add(time.Hour)}
 if err := Validate(c); err != nil { t.Fatalf("valid certificate rejected: %v", err) }
 if !IsExpiringSoon(c, now, 2*time.Hour) { t.Fatal("certificate inside expiry window not detected") }
 if IsExpiringSoon(c, now, -time.Hour) { t.Fatal("negative expiry window accepted") }
 if err := Validate(Certificate{}); err == nil { t.Fatal("empty certificate accepted") }
}
