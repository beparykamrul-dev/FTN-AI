package controlplane

import (
	"crypto/sha256"
	"encoding/hex"
)

// DeploymentEnvelope binds an approved plan to its exact content before it is
// handed to a transport layer. It contains no credentials or private keys.
type DeploymentEnvelope struct {
	Plan   DeploymentPlan
	Digest string
}

func SealPlan(p DeploymentPlan) DeploymentEnvelope {
	h := sha256.New()
	h.Write([]byte(p.ServerID))
	h.Write([]byte{0})
	for _, service := range p.Services {
		h.Write([]byte(service))
		h.Write([]byte{0})
	}
	if p.Approved { h.Write([]byte("approved")) }
	return DeploymentEnvelope{Plan: p, Digest: hex.EncodeToString(h.Sum(nil))}
}

// Verify checks that an envelope still represents the same approved plan.
func VerifyEnvelope(e DeploymentEnvelope) bool {
	return e.Plan.Approved && e.Digest == SealPlan(e.Plan).Digest
}
