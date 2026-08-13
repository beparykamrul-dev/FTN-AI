package proxy

import "time"

// CertificateAuthorityBoundary keeps certificate lifecycle authority outside
// request handlers. The proxy receives an opaque certificate version only;
// private key material must remain in the FTN certificate controller/HSM.
type CertificateAuthorityBoundary struct {
	RotationInterval time.Duration
	LastVersion      string
	LastRotated      time.Time
}

// CertificateDecision describes whether the control-plane certificate version
// may be activated. Request callers cannot select or inject private key data.
type CertificateDecision struct {
	Allowed bool
	Version string
	Reason  string
}

func (b *CertificateAuthorityBoundary) Activate(version string, now time.Time, securityPathHealthy bool) CertificateDecision {
	if version == "" { return CertificateDecision{Reason: "empty certificate version"} }
	if !securityPathHealthy { return CertificateDecision{Reason: "security path unhealthy"} }
	interval := b.RotationInterval
	if interval <= 0 { interval = time.Hour }
	if !b.LastRotated.IsZero() && now.Sub(b.LastRotated) < interval {
		return CertificateDecision{Version: b.LastVersion, Reason: "rotation interval not reached"}
	}
	b.LastVersion = version
	b.LastRotated = now
	return CertificateDecision{Allowed: true, Version: version, Reason: "certificate version activated"}
}
