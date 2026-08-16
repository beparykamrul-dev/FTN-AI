package edge

// ProviderTLS describes the verified TLS identity of an external edge/CDN/DNS provider.
type ProviderTLS struct {
	ProviderID  string
	Host        string
	Fingerprint string
	Verified    bool
	ExpiresUnix int64
}

func (p ProviderTLS) Valid(nowUnix int64) bool {
	return p.ProviderID != "" && p.Host != "" && p.Fingerprint != "" && p.Verified && p.ExpiresUnix > nowUnix
}
