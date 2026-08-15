package ftnservice

// DependencyEdge describes a service-to-service dependency.
type DependencyEdge struct {
	FromService string
	ToService   string
	Required    bool
}

// DependencyGraph is the declarative service dependency map.
type DependencyGraph struct {
	Edges []DependencyEdge
}

func (g DependencyGraph) Valid() bool {
	for _, e := range g.Edges {
		if e.FromService == "" || e.ToService == "" || e.FromService == e.ToService {
			return false
		}
	}
	return true
}
