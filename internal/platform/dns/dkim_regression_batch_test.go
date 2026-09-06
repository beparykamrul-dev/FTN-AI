package dns

import "testing"

func TestDKIMRejectsInvalidKeyType(t *testing.T) {
	r := DKIMRecord{Domain: "example.com", Selector: "default", PublicKey: "abc", KeyType: "unknown"}
	if r.Validate() == nil { t.Fatal("unsupported DKIM key type must be rejected") }
}

func TestDKIMFQDN(t *testing.T) {
	r := DKIMRecord{Domain: "Example.COM.", Selector: "mail."}
	if got := r.FQDN(); got != "mail._domainkey.Example.COM" { t.Fatalf("unexpected FQDN: %q", got) }
}
