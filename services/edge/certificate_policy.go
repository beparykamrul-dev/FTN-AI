package edge

// CertificatePolicy defines FTN certificate lifecycle requirements.
type CertificatePolicy struct {
	Provider string
	Domain   string
	KeyType  string
	AutoRenew bool
	MinDaysRemaining uint32
}

func (p CertificatePolicy) Valid() bool {
	return p.Provider != "" && p.Domain != "" && p.KeyType != "" && p.MinDaysRemaining > 0
}
