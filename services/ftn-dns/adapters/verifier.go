package adapters

import "context"

// Verifier is the common DNSSEC cryptographic verification contract.
// Implementations may use software crypto, HSMs, or external validation
// engines, but must return normalized results.
type Verifier interface {
	Name() string
	VerifyDS(ctx context.Context, dnskey []byte, dsDigest []byte, digestType uint8) error
	VerifyRRSIG(ctx context.Context, algorithm uint8, publicKey []byte, signedData []byte, signature []byte) error
}

// VerificationReport is consumed by FTN DNS trust and routing layers.
type VerificationReport struct {
	Provider      string
	DSVerified    bool
	RRSIGVerified bool
	ChainValid    bool
	Error         string
}

func (r VerificationReport) Trusted() bool {
	return r.Error == "" && r.DSVerified && r.RRSIGVerified && r.ChainValid
}
