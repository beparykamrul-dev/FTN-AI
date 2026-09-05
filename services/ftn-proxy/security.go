package proxy

import (
	"strings"
	"time"
)

// SecurityPolicy defines fail-closed limits for the FTN proxy data plane.
type SecurityPolicy struct {
	MaxConnectionsPerClient int
	MaxRequestsPerWindow    int
	RateWindow              time.Duration
	RequireTLS              bool
	RejectPrivateTargets    bool
	MaxHeaderBytes          int
}

func DefaultSecurityPolicy() SecurityPolicy {
	return SecurityPolicy{MaxConnectionsPerClient: 64, MaxRequestsPerWindow: 600, RateWindow: time.Minute, RequireTLS: true, RejectPrivateTargets: true, MaxHeaderBytes: 32 << 10}
}

// ValidateRequest applies cheap security checks before proxy forwarding.
// Target classification is deliberately supplied by the caller so this layer
// does not make assumptions about network topology.
func (p SecurityPolicy) ValidateRequest(scheme string, connections, requests int, headerBytes int, privateTarget bool) bool {
	scheme = strings.ToLower(strings.TrimSpace(scheme))
	if connections < 0 || requests < 0 || headerBytes < 0 {
		return false
	}
	if p.RequireTLS && scheme != "https" {
		return false
	}
	if p.MaxConnectionsPerClient > 0 && connections > p.MaxConnectionsPerClient {
		return false
	}
	if p.MaxRequestsPerWindow > 0 && requests > p.MaxRequestsPerWindow {
		return false
	}
	if p.MaxHeaderBytes > 0 && headerBytes > p.MaxHeaderBytes {
		return false
	}
	if p.RejectPrivateTargets && privateTarget {
		return false
	}
	return true
}
