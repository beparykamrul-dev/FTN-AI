package module

import "fmt"

// CapabilityConflicts reports capabilities claimed by more than one module.
// Duplicate capabilities are not automatically errors: the control plane can
// use this report to require an explicit provider/priority policy.
func (r *Registry) CapabilityConflicts() map[string][]string {
	owners := make(map[string][]string)
	for name, m := range r.modules {
		for _, capability := range m.Definition().Capabilities {
			owners[capability] = append(owners[capability], name)
		}
	}
	conflicts := make(map[string][]string)
	for capability, names := range owners {
		if len(names) > 1 { conflicts[capability] = names }
	}
	return conflicts
}

func (r *Registry) ValidateConflicts() error {
	for capability, names := range r.CapabilityConflicts() {
		return fmt.Errorf("capability %q has multiple providers: %v", capability, names)
	}
	return nil
}
