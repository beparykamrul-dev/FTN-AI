package module

import "fmt"

// ResolveOrder returns a dependency-first module order and rejects cycles.
func (r *Registry) ResolveOrder() ([]Definition, error) {
	state := make(map[string]uint8)
	order := make([]Definition, 0, len(r.modules))
	var visit func(string) error
	visit = func(name string) error {
		if state[name] == 1 {
			return fmt.Errorf("module dependency cycle at %q", name)
		}
		if state[name] == 2 {
			return nil
		}
		m, ok := r.modules[name]
		if !ok {
			return fmt.Errorf("missing module dependency %q", name)
		}
		state[name] = 1
		for _, dep := range m.Definition().Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[name] = 2
		order = append(order, m.Definition())
		return nil
	}
	for name := range r.modules {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	return order, nil
}
