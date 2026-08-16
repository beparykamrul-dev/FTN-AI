package edge

// HTTPSPolicy defines FTN Edge TLS requirements for external providers.
type HTTPSPolicy struct {
	Enabled          bool
	MinTLSVersion    string
	RequireSNI       bool
	RequireHostMatch bool
	AllowHTTPRedirect bool
}

func (p HTTPSPolicy) Valid() bool {
	return p.Enabled && (p.MinTLSVersion == "1.2" || p.MinTLSVersion == "1.3")
}
