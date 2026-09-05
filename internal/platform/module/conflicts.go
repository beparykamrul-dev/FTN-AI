package module

import (
	"fmt"
	"sort"
)

func (r *Registry) CapabilityConflicts() map[string][]string {
	conflicts := make(map[string][]string)
	if r == nil { return conflicts }
	r.mu.RLock()
	owners := make(map[string][]string)
	for name, m := range r.modules {
		if m == nil { continue }
		for _, capability := range m.Definition().Capabilities {
			if capability != "" { owners[capability] = append(owners[capability], name) }
		}
	}
	r.mu.RUnlock()
	for capability, names := range owners {
		if len(names) > 1 { sort.Strings(names); conflicts[capability] = names }
	}
	return conflicts
}

func (r *Registry) ValidateConflicts() error {
	conflicts := r.CapabilityConflicts()
	keys := make([]string, 0, len(conflicts))
	for capability := range conflicts { keys = append(keys, capability) }
	sort.Strings(keys)
	if len(keys) == 0 { return nil }
	capability := keys[0]
	return fmt.Errorf("capability %q has multiple providers: %v", capability, conflicts[capability])
}
