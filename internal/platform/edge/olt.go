package edge

import (
	"context"
	"strings"
)

type OLTState struct {
	ID           string `json:"id"`
	Vendor       string `json:"vendor"`
	Model        string `json:"model,omitempty"`
	ManagementIP string `json:"management_ip"`
	Status       string `json:"status"`
}

type ONUState struct {
	ID         string  `json:"id"`
	Serial     string  `json:"serial"`
	PON        string  `json:"pon"`
	Status     string  `json:"status"`
	RXPowerDBM float64 `json:"rx_power_dbm,omitempty"`
	TXPowerDBM float64 `json:"tx_power_dbm,omitempty"`
}

type OLTDriver interface {
	Identity(ctx context.Context) (OLTState, error)
	DiscoverONUs(ctx context.Context) ([]ONUState, error)
	Apply(ctx context.Context, request ChangeRequest) error
}

func (s OLTState) Valid() bool {
	return strings.TrimSpace(s.ID) != "" && strings.TrimSpace(s.Vendor) != "" && strings.TrimSpace(s.ManagementIP) != ""
}

func (s ONUState) Valid() bool {
	return strings.TrimSpace(s.ID) != "" && strings.TrimSpace(s.Serial) != "" && strings.TrimSpace(s.PON) != ""
}
