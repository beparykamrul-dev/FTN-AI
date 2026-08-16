package edge

// ProviderTLSVerifier describes the normalized result of validating an
// external CDN/Edge/DNS provider TLS connection.
type ProviderTLSVerifier struct {
	Provider   string
	Domain     string
	SNI        string
	TLSVersion string
	Verified   bool
	ErrorClass string
}

func (v ProviderTLSVerifier) Valid(policy HTTPSPolicy) bool {
	return v.Verified && v.Provider != "" && v.Domain != "" &&
		v.SNI == v.Domain && (v.TLSVersion == "1.2" || v.TLSVersion == "1.3") && policy.Valid()
}
