package controlplane

import (
	"crypto/sha256"
	"encoding/hex"
)

// AgentIdentity is the non-secret identity presented by an enrolled FTN server.
// Private keys remain on the server agent or approved identity infrastructure.
type AgentIdentity struct {
	ServerID    string
	Fingerprint string
	Enrolled    bool
	Revoked     bool
}

// AuthorizeAgent accepts only enrolled, non-revoked identities with a stable
// fingerprint. Transport authentication (for example mTLS) is performed by
// the network layer; this function is the control-plane authorization gate.
func AuthorizeAgent(id AgentIdentity, presentedFingerprint string) bool {
	return id.Enrolled && !id.Revoked && id.Fingerprint != "" &&
		presentedFingerprint != "" && id.Fingerprint == presentedFingerprint
}

func Fingerprint(identityMaterial []byte) string {
	sum := sha256.Sum256(identityMaterial)
	return hex.EncodeToString(sum[:])
}
