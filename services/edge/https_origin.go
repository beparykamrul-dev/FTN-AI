package edge

// HTTPSOriginPolicy controls secure forwarding from FTN Edge to an origin.
type HTTPSOriginPolicy struct {
	OriginHost      string
	RequireTLS      bool
	VerifyCertificate bool
	AllowedSNI      string
}

func (p HTTPSOriginPolicy) Valid() bool {
	return p.OriginHost != "" && p.RequireTLS && p.VerifyCertificate
}
