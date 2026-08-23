package module

import (
	"sort"
	"sync"
)

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
	mu      sync.RWMutex
	modules map[string]Module
}

func NewRegistry() *Registry { return &Registry{modules: make(map[string]Module)} }

func (r *Registry) Register(m Module) {
	if m == nil {
		return
	}
	d := m.Definition()
	if d.Name == "" {
		return
	}
	r.mu.Lock()
	r.modules[d.Name] = m
	r.mu.Unlock()
}

func (r *Registry) Get(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.modules[name]
	return m, ok
}

func (r *Registry) List() []Definition {
	r.mu.RLock()
	out := make([]Definition, 0, len(r.modules))
	for _, m := range r.modules {
		out = append(out, m.Definition())
	}
	r.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.modules[name]
	return ok
}

// DependenciesReady verifies that every declared module dependency is loaded.
func (r *Registry) DependenciesReady(d Definition) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, dep := range d.Dependencies {
		if _, ok := r.modules[dep]; !ok {
			return false
		}
	}
	return true
}
