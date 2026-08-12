package transport

import (
	"errors"
	"net"
	"strings"
)

var (
	ErrRelayDenied = errors.New("relay denied")
	ErrTLSRequired = errors.New("tls required")
	ErrAuthRequired = errors.New("authentication required")
)

type SMTPPolicy struct {
	LocalDomains map[string]bool
	RequireTLS   bool
	RequireAuth  bool
}

func (p SMTPPolicy) ValidateSubmission(remote net.IP, tlsActive, authenticated bool, from, recipient string) error {
	if p.RequireTLS && !tlsActive { return ErrTLSRequired }
	if p.RequireAuth && !authenticated { return ErrAuthRequired }
	if from == "" || recipient == "" { return errors.New("sender and recipient required") }
	parts := strings.SplitN(recipient, "@", 2)
	if len(parts) != 2 || parts[0] == "" || !p.LocalDomains[strings.ToLower(parts[1])] {
		return ErrRelayDenied
	}
	_ = remote // reserved for rate-limit / abuse policy at the connection layer
	return nil
}
