package controlplane

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

// AgentIdentity is the public identity presented by an enrolled FTN agent.
type AgentIdentity struct {
	ServerID    string
	Fingerprint string
	Enrolled    bool
	Revoked     bool
}

// AuthorizeAgent provides a constant-time comparison for the identity
// fingerprint. Cryptographic proof of possession remains the responsibility
// of the mTLS transport layer.
func AuthorizeAgent(expected, presented AgentIdentity) bool {
	if expected.ServerID == "" || presented.ServerID == "" || expected.ServerID != presented.ServerID {
		return false
	}
	if !expected.Enrolled || !presented.Enrolled || expected.Revoked || presented.Revoked || expected.Fingerprint == "" || presented.Fingerprint == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected.Fingerprint), []byte(presented.Fingerprint)) == 1
}

// Fingerprint derives a stable public identity fingerprint from identity material.
// Private key material must never be passed here.
func Fingerprint(identityMaterial []byte) string {
	sum := sha256.Sum256(identityMaterial)
	return hex.EncodeToString(sum[:])
}
