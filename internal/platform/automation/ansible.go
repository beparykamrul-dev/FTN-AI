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
	for _, host := range p.Limit { if strings.TrimSpace(host) == "" { return fmt.Errorf("limit contains an empty target") } }
	return nil
}

func (p AnsiblePlan) Normalize() AnsiblePlan {
	p.Inventory = strings.TrimSpace(p.Inventory)
	p.Playbook = strings.TrimSpace(p.Playbook)
	seen := make(map[string]struct{}, len(p.Limit))
	out := make([]string, 0, len(p.Limit))
	for _, host := range p.Limit { host = strings.TrimSpace(host); if host == "" { continue }; if _, ok := seen[host]; ok { continue }; seen[host] = struct{}{}; out = append(out, host) }
	p.Limit = out
	return p
}
