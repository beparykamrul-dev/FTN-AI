package agent

import "sort"

// Capability describes a task a layer can perform.
type Capability struct {
	Name     string
	Category Category
	Priority int
}

// MatchCapabilities orders capabilities for a requested category without binding
// FTN's API to a particular model/provider.
func MatchCapabilities(category Category, capabilities []Capability) []Capability {
	matched := make([]Capability, 0, len(capabilities))
	for _, c := range capabilities {
		if c.Category == category {
			matched = append(matched, c)
		}
	}
	sort.SliceStable(matched, func(i, j int) bool {
		return matched[i].Priority < matched[j].Priority
	})
	return matched
}
