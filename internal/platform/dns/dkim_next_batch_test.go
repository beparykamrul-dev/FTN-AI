package dns

import "testing"

func TestDKIMValidationAndFQDN(t *testing.T) {
	r := DKIMRecord{Domain: " Example.COM. ", Selector: "mail.", PublicKey: "abc", KeyType: "rsa", Hash: "sha256"}
	if err := r.Validate(); err != nil { t.Fatalf("valid DKIM rejected: %v", err) }
	if got := r.FQDN(); got != "mail._domainkey.Example.COM" { t.Fatalf("FQDN=%q", got) }
	if got, err := r.TXTValue(); err != nil || got != "v=DKIM1; k=rsa; p=abc; h=sha256" { t.Fatalf("TXT=%q err=%v", got, err) }
	if err := (DKIMRecord{Domain: "example.com", Selector: "s", PublicKey: "k", KeyType: "ed25519", Hash: "sha1"}).Validate(); err == nil { t.Fatal("unsupported DKIM hash accepted") }
}
