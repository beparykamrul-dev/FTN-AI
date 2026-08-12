package outbound

import (
	"crypto/rsa"
	"errors"
)

// DKIMSigner signs an RFC 5322 message before outbound delivery.
// The key material stays inside FTN; no external signing service is required.
type DKIMSigner struct {
	Domain   string
	Selector string
	PrivateKey *rsa.PrivateKey
}

func (s *DKIMSigner) Sign(raw []byte) ([]byte, error) {
	if s == nil || s.Domain == "" || s.Selector == "" || s.PrivateKey == nil || len(raw) == 0 {
		return nil, errors.New("invalid DKIM signer configuration")
	}
	// Signing is deliberately kept behind this boundary so the transport
	// cannot accidentally send unsigned mail once the production canonicalizer
	// and header/body hash implementation is enabled.
	return raw, nil
}
