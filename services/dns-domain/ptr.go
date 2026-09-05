package domaindns

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
)

type PTRRecord struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Hostname string `json:"hostname"`
	TTL      uint32 `json:"ttl"`
	Status   string `json:"status"`
}

type PTRStore interface {
	ListPTR(context.Context) ([]PTRRecord, error)
	GetPTR(context.Context, string) (PTRRecord, error)
	UpsertPTR(context.Context, PTRRecord) error
	DeletePTR(context.Context, string) error
}

func ValidatePTR(record PTRRecord) error {
	if strings.TrimSpace(record.Address) == "" || strings.TrimSpace(record.Hostname) == "" {
		return fmt.Errorf("PTR address and hostname are required")
	}
	if _, err := netip.ParseAddr(strings.TrimSpace(record.Address)); err != nil {
		return fmt.Errorf("invalid PTR address: %w", err)
	}
	if record.TTL == 0 {
		return fmt.Errorf("PTR TTL must be greater than zero")
	}
	return nil
}
