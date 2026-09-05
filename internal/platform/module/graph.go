package module

import (
	"fmt"
	"sort"
	"strings"
)

func (r *Registry) ResolveOrder() ([]Definition, error) {
	if r == nil { return nil, fmt.Errorf("module registry is required") }
	r.mu.RLock()
	defer r.mu.RUnlock()
	state := make(map[string]uint8)
	order := make([]Definition, 0, len(r.modules))
	var visit func(string) error
	visit = func(name string) error {
		name = strings.TrimSpace(name)
		if name == "" { return fmt.Errorf("module name is required") }
		if state[name] == 1 { return fmt.Errorf("module dependency cycle at %q", name) }
		if state[name] == 2 { return nil }
		m, ok := r.modules[name]
		if !ok || m == nil { return fmt.Errorf("missing module dependency %q", name) }
		state[name] = 1
		deps := append([]string(nil), m.Definition().Dependencies...)
		sort.Strings(deps)
		seen := make(map[string]struct{}, len(deps))
		for _, dep := range deps {
			dep = strings.TrimSpace(dep)
			if dep == "" { continue }
			if _, ok := seen[dep]; ok { continue }
			seen[dep] = struct{}{}
			if err := visit(dep); err != nil { return err }
		}
		state[name] = 2
		order = append(order, m.Definition())
		return nil
	}
	names := make([]string, 0, len(r.modules))
	for name := range r.modules { names = append(names, name) }
	sort.Strings(names)
	for _, name := range names { if err := visit(name); err != nil { return nil, err } }
	return order, nil
}
