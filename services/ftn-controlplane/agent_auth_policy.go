package controlplane

// AgentAuthPolicy defines the minimum trust requirements for an FTN server
// agent. TLS proof and certificate validation remain in the transport layer.
type AgentAuthPolicy struct {
	RequireMTLS       bool
	RequireEnrollment bool
	RejectRevoked     bool
}

func DefaultAgentAuthPolicy() AgentAuthPolicy {
	return AgentAuthPolicy{RequireMTLS: true, RequireEnrollment: true, RejectRevoked: true}
}

// AuthorizePeer applies the control-plane policy to an already TLS-verified
// peer. It never accepts private keys or performs cryptographic handshakes.
func (p AgentAuthPolicy) AuthorizePeer(expected AgentIdentity, peer MTLSIdentity, enrolled bool) bool {
	if p.RequireMTLS && !peer.Verified { return false }
	if p.RequireEnrollment && !enrolled { return false }
	if p.RejectRevoked && peer.Revoked { return false }
	return AuthorizeAgent(expected, AgentIdentity{ServerID: peer.ServerID, Fingerprint: peer.Fingerprint})
}
