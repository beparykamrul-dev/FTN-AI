package edge

// TLSHealth captures the normalized certificate/TLS health state for an edge endpoint.
type TLSHealth struct {
	Domain string
	Valid  bool
	DaysRemaining uint32
	TLSVersion string
	ErrorClass string
}

func (h TLSHealth) Healthy(policy CertificatePolicy) bool {
	return h.Valid && h.Domain == policy.Domain && h.DaysRemaining >= policy.MinDaysRemaining
}
