package graph

import (
	"sort"
	"strings"
	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

// Graph provides a deterministic dependency view for authorized service relationships.
type Graph struct { Dependencies []model.Dependency }

func (g Graph) DependenciesOf(service string) []model.Dependency {
	service = strings.TrimSpace(service); if service == "" { return nil }
	out := make([]model.Dependency, 0)
	seen := make(map[string]struct{})
	for _, dep := range g.Dependencies { dep = dep.Normalize(); if !dep.Valid() || (dep.From != service && dep.To != service) { continue }; key := dep.From+"\x00"+dep.To+"\x00"+dep.Kind; if _, ok := seen[key]; ok { continue }; seen[key] = struct{}{}; out = append(out, dep) }
	sort.SliceStable(out, func(i,j int) bool { if out[i].From != out[j].From { return out[i].From < out[j].From }; if out[i].To != out[j].To { return out[i].To < out[j].To }; return out[i].Kind < out[j].Kind })
	return out
}
