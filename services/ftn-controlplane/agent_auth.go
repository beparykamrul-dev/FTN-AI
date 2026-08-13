package controlplane

import "crypto/subtle"

// AgentIdentity is the public identity presented by an enrolled FTN agent.
type AgentIdentity struct {
	ServerID   string
	Fingerprint string
	Revoked    bool
}

// AuthorizeAgent provides a constant-time comparison for the identity
// fingerprint. Cryptographic proof of possession remains the responsibility
// of the mTLS transport layer.
func AuthorizeAgent(expected, presented AgentIdentity) bool {
	if expected.ServerID == "" || presented.ServerID == "" || expected.ServerID != presented.ServerID {
		return false
	}
	if expected.Revoked || presented.Revoked || expected.Fingerprint == "" || presented.Fingerprint == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected.Fingerprint), []byte(presented.Fingerprint)) == 1
}
