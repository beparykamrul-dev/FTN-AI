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

func validPTRHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 253 { return false }
	for _, label := range strings.Split(hostname, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' { return false }
		for _, r := range label {
			if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') && r != '-' { return false }
		}
	}
	return true
}

func ValidatePTR(record PTRRecord) error {
	if strings.TrimSpace(record.ID) == "" { return fmt.Errorf("PTR id is required") }
	address := strings.TrimSpace(record.Address)
	hostname := strings.TrimSuffix(strings.TrimSpace(record.Hostname), ".")
	if address == "" || hostname == "" { return fmt.Errorf("PTR address and hostname are required") }
	if _, err := netip.ParseAddr(address); err != nil { return fmt.Errorf("invalid PTR address: %w", err) }
	if !validPTRHostname(hostname) { return fmt.Errorf("invalid PTR hostname") }
	if record.TTL == 0 { return fmt.Errorf("PTR TTL must be greater than zero") }
	return nil
}
