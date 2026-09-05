package dns

import "testing"

func TestDKIMRejectsUnsupportedKeyType(t *testing.T) {
	r := DKIMRecord{Domain:"example.com",Selector:"default",PublicKey:"abc",KeyType:"bad"}
	if err := r.Validate(); err == nil { t.Fatal("unsupported DKIM key type must be rejected") }
}

func TestDKIMFQDNIsCanonical(t *testing.T) {
	r := DKIMRecord{Domain:" example.com. ",Selector:" selector. ",PublicKey:"abc"}
	if got := r.FQDN(); got != "selector._domainkey.example.com" { t.Fatalf("unexpected FQDN: %q", got) }
}
