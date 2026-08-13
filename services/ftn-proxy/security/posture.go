package security

// Posture describes security capabilities available to FTN services without
// coupling the proxy to a particular security product.
type Posture struct {
	TLS             bool
	Identity        bool
	ReplayProtection bool
	Audit           bool
	RateLimit       bool
	NetworkPolicy   bool
	KeyIsolation    bool
}

// ReadyForSecurityService reports whether the infrastructure exposes the
// control points required to add a dedicated cyber-security service later.
func (p Posture) ReadyForSecurityService() bool {
	return p.TLS && p.Identity && p.ReplayProtection && p.Audit &&
		p.RateLimit && p.NetworkPolicy && p.KeyIsolation
}
