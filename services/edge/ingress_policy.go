package edge

// IngressPolicy defines how FTN accepts traffic from approved external
// providers before forwarding it into FTN Edge/FTNWAN.
type IngressPolicy struct {
	Provider string
	RequireTLS bool
	RequireVerifiedSource bool
	AllowedHosts []string
}

func (p IngressPolicy) Allows(provider, host string, tls, verified bool) bool {
	if p.Provider != provider || (p.RequireTLS && !tls) || (p.RequireVerifiedSource && !verified) {
		return false
	}
	for _, h := range p.AllowedHosts {
		if h == host { return true }
	}
	return false
}
