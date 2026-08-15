package ftnservice

// Service describes one FTN service in the unified control plane.
type Service struct {
	ID       string
	Version  string
	Region   string
	Healthy  bool
	Enabled  bool
}

// Registry is the logical FTN service registry.
type Registry struct {
	Services map[string]Service
}

func (r Registry) Get(id string) (Service, bool) {
	s, ok := r.Services[id]
	return s, ok
}
