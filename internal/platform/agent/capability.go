package agent

import (
	"sort"
	"strings"
)

// Capability describes a task a layer can perform.
type Capability struct {
	Name     string
	Category Category
	Priority int
}

// MatchCapabilities returns only usable capabilities in deterministic order.
func MatchCapabilities(category Category, capabilities []Capability) []Capability {
	matched := make([]Capability, 0, len(capabilities))
	for _, c := range capabilities {
		c.Name = strings.TrimSpace(c.Name)
		if c.Name == "" || c.Category != category {
			continue
		}
		matched = append(matched, c)
	}
	sort.SliceStable(matched, func(i, j int) bool {
		if matched[i].Priority != matched[j].Priority {
			return matched[i].Priority < matched[j].Priority
		}
		return matched[i].Name < matched[j].Name
	})
	return matched
}
