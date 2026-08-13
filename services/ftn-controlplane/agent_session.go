package controlplane

import "time"

// AgentSession represents a short-lived, authenticated control-plane session.
// Session state contains no private key material or customer payloads.
type AgentSession struct {
	ServerID    string
	Fingerprint string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	Revoked     bool
}

func (s AgentSession) Valid(now time.Time) bool {
	if s.ServerID == "" || s.Fingerprint == "" || s.Revoked {
		return false
	}
	return !now.Before(s.IssuedAt) && now.Before(s.ExpiresAt)
}
