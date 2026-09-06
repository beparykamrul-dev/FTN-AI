package domaindns

import "testing"

func TestPTRValidationBoundaries(t *testing.T) {
 r := PTRRecord{ID: "ptr-1", Address: "192.0.2.10", Hostname: "host.example.com.", TTL: 300}
 if err := ValidatePTR(r); err != nil { t.Fatalf("valid PTR rejected: %v", err) }
 r.Address = "not-an-ip"; if err := ValidatePTR(r); err == nil { t.Fatal("invalid PTR address accepted") }
 r = PTRRecord{ID: "ptr-1", Address: "192.0.2.10", Hostname: "-bad.example.com", TTL: 300}; if err := ValidatePTR(r); err == nil { t.Fatal("invalid PTR hostname accepted") }
}
