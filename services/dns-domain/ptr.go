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
	address := strings.TrimSpace(record.Address)
	hostname := strings.TrimSuffix(strings.TrimSpace(record.Hostname), ".")
	if address == "" || hostname == "" {
		return fmt.Errorf("PTR address and hostname are required")
	}
	if _, err := netip.ParseAddr(address); err != nil {
		return fmt.Errorf("invalid PTR address: %w", err)
	}
	if strings.ContainsAny(hostname, " \t\r\n") {
		return fmt.Errorf("invalid PTR hostname")
	}
	if record.TTL == 0 {
		return fmt.Errorf("PTR TTL must be greater than zero")
	}
	return nil
}
