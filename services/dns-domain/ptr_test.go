package domaindns

import "testing"

func TestValidatePTRRejectsInvalidInput(t *testing.T) {
	cases := []PTRRecord{
		{ID: "", Address: "192.0.2.1", Hostname: "host.example.com", TTL: 300},
		{ID: "p1", Address: "not-an-ip", Hostname: "host.example.com", TTL: 300},
		{ID: "p1", Address: "192.0.2.1", Hostname: "-bad.example.com", TTL: 300},
		{ID: "p1", Address: "192.0.2.1", Hostname: "host.example.com", TTL: 604801},
	}
	for _, tc := range cases {
		if err := ValidatePTR(tc); err == nil { t.Fatalf("expected invalid PTR to fail: %+v", tc) }
	}
}

func TestValidatePTRAcceptsTrailingDot(t *testing.T) {
	if err := ValidatePTR(PTRRecord{ID: "p1", Address: "2001:db8::1", Hostname: "host.example.com.", TTL: 300}); err != nil { t.Fatalf("expected valid PTR: %v", err) }
}
