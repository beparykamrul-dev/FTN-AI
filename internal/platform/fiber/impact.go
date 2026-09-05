package fiber

import "sort"

// ImpactNode describes an infrastructure/customer dependency reached by a
// fiber path. Customer identifiers are references only; authorization belongs
// to the control-plane policy layer.
type ImpactNode struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
}

type FiberPath struct {
	ID       string       `json:"id"`
	Segments []string     `json:"segments"`
	Impacted []ImpactNode `json:"impacted"`
}

// BuildImpact returns a deterministic, de-duplicated impact set for a fiber
// path. It does not infer outages; callers must provide observed impact data.
func BuildImpact(path FiberPath) []ImpactNode {
	seen := make(map[string]ImpactNode, len(path.Impacted))
	for _, n := range path.Impacted {
		if n.ID != "" {
			seen[n.ID] = n
		}
	}
	out := make([]ImpactNode, 0, len(seen))
	for _, n := range seen {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ID != out[j].ID {
			return out[i].ID < out[j].ID
		}
		return out[i].Kind < out[j].Kind
	})
	return out
}
