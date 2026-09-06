package automation

import (
	"fmt"
	"sort"
	"strings"
)

type AnsiblePlan struct {
	Inventory string   `json:"inventory"`
	Playbook  string   `json:"playbook"`
	Limit     []string `json:"limit,omitempty"`
	CheckMode bool     `json:"check_mode"`
}

func (p AnsiblePlan) Validate() error {
	inventory := strings.TrimSpace(p.Inventory)
	playbook := strings.TrimSpace(p.Playbook)
	if inventory == "" {
		return fmt.Errorf("inventory is required")
	}
	if playbook == "" {
		return fmt.Errorf("playbook is required")
	}
	if len(inventory) > 4096 || len(playbook) > 4096 {
		return fmt.Errorf("inventory or playbook is too long")
	}
	if strings.ContainsAny(inventory, "\r\n\x00") || strings.ContainsAny(playbook, "\r\n\x00") {
		return fmt.Errorf("inventory or playbook contains control characters")
	}
	if len(p.Limit) > 1024 {
		return fmt.Errorf("too many limit targets")
	}
	for _, host := range p.Limit {
		host = strings.TrimSpace(host)
		if host == "" || len(host) > 512 || strings.ContainsAny(host, "\r\n\x00") {
			return fmt.Errorf("invalid limit target")
		}
	}
	return nil
}

func (p AnsiblePlan) Normalize() AnsiblePlan {
	p.Inventory = strings.TrimSpace(p.Inventory)
	p.Playbook = strings.TrimSpace(p.Playbook)
	seen := map[string]struct{}{}
	out := make([]string, 0, len(p.Limit))
	for _, host := range p.Limit {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}
		if _, ok := seen[host]; ok {
			continue
		}
		seen[host] = struct{}{}
		out = append(out, host)
	}
	sort.Strings(out)
	p.Limit = out
	return p
}
