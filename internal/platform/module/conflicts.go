package module

import (
	"fmt"
	"sort"
)

// CapabilityConflicts reports capabilities claimed by more than one module.
// Duplicate capabilities are not automatically errors: the control plane can
// use this report to require an explicit provider/priority policy.
func (r *Registry) CapabilityConflicts() map[string][]string {
	conflicts := make(map[string][]string)
	if r == nil { return conflicts }
	r.mu.RLock()
	owners := make(map[string][]string)
	for name, m := range r.modules {
		if m == nil { continue }
		for _, capability := range m.Definition().Capabilities { if capability != "" { owners[capability] = append(owners[capability], name) } }
	}
	r.mu.RUnlock()
	for capability, names := range owners { if len(names)>1 { sort.Strings(names); conflicts[capability]=names } }
	return conflicts
}
func (r *Registry) ValidateConflicts() error { for capability,names:=range r.CapabilityConflicts(){return fmt.Errorf("capability %q has multiple providers: %v",capability,names)};return nil }
