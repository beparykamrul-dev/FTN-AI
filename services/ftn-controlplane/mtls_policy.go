package controlplane

type MTLSIdentity struct {
	ServerID string
	Fingerprint string
	Verified bool
	Revoked bool
}

func VerifyMTLSIdentity(expected AgentIdentity, peer MTLSIdentity) bool {
	if !peer.Verified || peer.Revoked { return false }
	return AuthorizeAgent(expected, AgentIdentity{ServerID: peer.ServerID, Fingerprint: peer.Fingerprint, Revoked: peer.Revoked})
}
