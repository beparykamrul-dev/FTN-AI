package automation

import (
	"fmt"
	"strings"
)

type AnsiblePlan struct {
	Inventory string   `json:"inventory"`
	Playbook  string   `json:"playbook"`
	Limit     []string `json:"limit,omitempty"`
	CheckMode bool     `json:"check_mode"`
}

func (p AnsiblePlan) Validate() error {
	if strings.TrimSpace(p.Inventory) == "" { return fmt.Errorf("inventory is required") }
	if strings.TrimSpace(p.Playbook) == "" { return fmt.Errorf("playbook is required") }
	return nil
}
