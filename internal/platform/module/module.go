package module

type Definition struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type Module interface {
	Definition() Definition
}

type Registry struct {
	modules map[string]Module
}

func NewRegistry() *Registry { return &Registry{modules: make(map[string]Module)} }

func (r *Registry) Register(m Module) { r.modules[m.Definition().Name] = m }
func (r *Registry) Get(name string) (Module, bool) { m, ok := r.modules[name]; return m, ok }
func (r *Registry) List() []Definition {
	out := make([]Definition, 0, len(r.modules))
	for _, m := range r.modules { out = append(out, m.Definition()) }
	return out
}

func (r *Registry) Has(name string) bool { _, ok := r.modules[name]; return ok }

// DependenciesReady verifies that every declared module dependency is loaded.
func (r *Registry) DependenciesReady(d Definition) bool {
	for _, dep := range d.Dependencies { if !r.Has(dep) { return false } }
	return true
}
