package module

import "fmt"

type Manifest struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// Validate checks the module manifest before runtime registration.
func (m Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("module name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("module version is required")
	}
	seen := make(map[string]struct{}, len(m.Capabilities))
	for _, c := range m.Capabilities {
		if c == "" {
			return fmt.Errorf("module %q contains an empty capability", m.Name)
		}
		if _, ok := seen[c]; ok {
			return fmt.Errorf("module %q declares duplicate capability %q", m.Name, c)
		}
		seen[c] = struct{}{}
	}
	for _, d := range m.Dependencies {
		if d == "" {
			return fmt.Errorf("module %q contains an empty dependency", m.Name)
		}
		if d == m.Name {
			return fmt.Errorf("module %q cannot depend on itself", m.Name)
		}
	}
	return nil
}
