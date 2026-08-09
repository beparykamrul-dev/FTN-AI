package inventory

import (
	"fmt"
	"strings"
)

type NetBoxObject struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Address  string            `json:"address,omitempty"`
	Site     string            `json:"site,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type NetBoxClient struct { BaseURL string; TokenRef string }

func (c NetBoxClient) Validate() error {
	if strings.TrimSpace(c.BaseURL) == "" { return fmt.Errorf("NetBox base URL is required") }
	if strings.TrimSpace(c.TokenRef) == "" { return fmt.Errorf("NetBox token reference is required") }
	return nil
}
