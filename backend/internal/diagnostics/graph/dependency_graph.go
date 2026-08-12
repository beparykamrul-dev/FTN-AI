package graph

import "github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"

// Graph provides a deterministic dependency view for authorized service relationships.
type Graph struct {
	Dependencies []model.Dependency
}

func (g Graph) DependenciesOf(service string) []model.Dependency {
	out := make([]model.Dependency, 0)
	for _, dep := range g.Dependencies {
		if dep.From == service || dep.To == service {
			out = append(out, dep)
		}
	}
	return out
}
