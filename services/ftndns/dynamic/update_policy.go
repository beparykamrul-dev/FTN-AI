package dynamic

import (
	"errors"
	"net"
	"strings"
	"time"
)

// UpdateRequest describes an authenticated FTNDNS dynamic update.
type UpdateRequest struct {
	Name      string
	Type      string
	Address   string
	NodeID    string
	TTL       uint32
	ExpiresAt time.Time
}

var ErrInvalidUpdate = errors.New("invalid FTNDNS update")

// ValidateUpdate validates an update before it reaches an authoritative DNS
// backend. Backend mutation and zone transfer remain outside this package.
func ValidateUpdate(r UpdateRequest) error {
	if strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.NodeID) == "" {
		return ErrInvalidUpdate
	}
	if r.Type != "A" && r.Type != "AAAA" {
		return ErrInvalidUpdate
	}
	if net.ParseIP(r.Address) == nil || r.TTL == 0 || r.TTL > 86400 {
		return ErrInvalidUpdate
	}
	if r.ExpiresAt.Before(time.Now().UTC()) {
		return ErrInvalidUpdate
	}
	return nil
}
