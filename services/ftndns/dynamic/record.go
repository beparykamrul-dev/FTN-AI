package dynamic

import (
	"errors"
	"net"
	"strings"
	"time"
)

// Record is an FTNDNS dynamic record managed by the FTN control plane.
type Record struct {
	Name      string
	Type      string
	Value     string
	TTL       uint32
	NodeID    string
	UpdatedAt time.Time
}

var ErrInvalidRecord = errors.New("invalid dynamic DNS record")

func ValidateRecord(r Record) error {
	if strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.NodeID) == "" || r.TTL == 0 {
		return ErrInvalidRecord
	}
	switch strings.ToUpper(r.Type) {
	case "A":
		if net.ParseIP(r.Value) == nil || net.ParseIP(r.Value).To4() == nil {
			return ErrInvalidRecord
		}
	case "AAAA":
		if net.ParseIP(r.Value) == nil || net.ParseIP(r.Value).To4() != nil {
			return ErrInvalidRecord
		}
	default:
		return ErrInvalidRecord
	}
	return nil
}
