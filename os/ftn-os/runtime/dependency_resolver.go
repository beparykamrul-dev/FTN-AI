package runtime

import "fmt"

// ResolveStartupOrder returns a deterministic topological startup order.
// Cycles and missing dependencies are rejected before any service starts.
func ResolveStartupOrder(modules []Module) ([]Module, error) {
	byName := make(map[string]Module, len(modules))
	for _, module := range modules {
		if module.Name == "" {
			return nil, fmt.Errorf("module name is required")
		}
		if _, exists := byName[module.Name]; exists {
			return nil, fmt.Errorf("duplicate module: %s", module.Name)
		}
		byName[module.Name] = module
	}

	state := make(map[string]uint8, len(modules))
	order := make([]Module, 0, len(modules))

	var visit func(string) error
	visit = func(name string) error {
		switch state[name] {
		case 1:
			return fmt.Errorf("dependency cycle detected at module: %s", name)
		case 2:
			return nil
		}

		module, ok := byName[name]
		if !ok {
			return fmt.Errorf("missing dependency: %s", name)
		}

		state[name] = 1
		for _, dependency := range module.Dependencies {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		state[name] = 2
		order = append(order, module)
		return nil
	}

	for _, module := range modules {
		if err := visit(module.Name); err != nil {
			return nil, err
		}
	}
	return order, nil
}
