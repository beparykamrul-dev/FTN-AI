package ftnwan

// OverlayKind identifies an FTNWAN connectivity mechanism.
type OverlayKind string

const (
	OverlayDirect       OverlayKind = "direct"
	OverlayWireGuard    OverlayKind = "wireguard"
	OverlayMTLS         OverlayKind = "mtls"
	OverlayProxy        OverlayKind = "proxy"
	OverlaySOCKS5       OverlayKind = "socks5"
	OverlayShadowsocks  OverlayKind = "shadowsocks"
	OverlayHeadscale    OverlayKind = "headscale"
	OverlayNetBird      OverlayKind = "netbird"
	OverlayZeroTier     OverlayKind = "zerotier"
	OverlayNebula       OverlayKind = "nebula"
)

// OverlayPath represents a discovered path; it contains no credentials or secrets.
type OverlayPath struct {
	Name      string
	Kind      OverlayKind
	Source    string
	Target    string
	Encrypted bool
	Healthy   bool
}

// SecureCandidate reports whether a path is eligible for policy evaluation.
func SecureCandidate(p OverlayPath) bool {
	if !p.Healthy {
		return false
	}
	if p.Kind == OverlayDirect {
		return true
	}
	return p.Encrypted
}
