package module

import (
	"sort"
	"strings"
	"sync"
)

type Definition struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type Module interface{ Definition() Definition }

type Registry struct {
	mu      sync.RWMutex
	modules map[string]Module
}

func NewRegistry() *Registry { return &Registry{modules: make(map[string]Module)} }

func normalizeDefinition(d Definition) Definition {
	d.Name = strings.TrimSpace(d.Name)
	d.Version = strings.TrimSpace(d.Version)
	caps := make([]string, 0, len(d.Capabilities))
	seenCaps := make(map[string]struct{}, len(d.Capabilities))
	for _, v := range d.Capabilities {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seenCaps[v]; ok {
			continue
		}
		seenCaps[v] = struct{}{}
		caps = append(caps, v)
	}
	sort.Strings(caps)
	d.Capabilities = caps
	deps := make([]string, 0, len(d.Dependencies))
	seenDeps := make(map[string]struct{}, len(d.Dependencies))
	for _, v := range d.Dependencies {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seenDeps[v]; ok {
			continue
		}
		seenDeps[v] = struct{}{}
		deps = append(deps, v)
	}
	sort.Strings(deps)
	d.Dependencies = deps
	return d
}

func (r *Registry) Register(m Module) {
	if r == nil || m == nil {
		return
	}
	d := normalizeDefinition(m.Definition())
	if d.Name == "" {
		return
	}
	r.mu.Lock()
	if r.modules == nil {
		r.modules = make(map[string]Module)
	}
	r.modules[d.Name] = m
	r.mu.Unlock()
}

func (r *Registry) Get(name string) (Module, bool) {
	if r == nil {
		return nil, false
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.modules[name]
	return m, ok
}

func (r *Registry) List() []Definition {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	out := make([]Definition, 0, len(r.modules))
	for _, m := range r.modules {
		if m != nil {
			out = append(out, normalizeDefinition(m.Definition()))
		}
	}
	r.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (r *Registry) Has(name string) bool {
	if r == nil {
		return false
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.modules[name]
	return ok
}

func (r *Registry) DependenciesReady(d Definition) bool {
	if r == nil {
		return false
	}
	d = normalizeDefinition(d)
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, dep := range d.Dependencies {
		if _, ok := r.modules[dep]; !ok {
			return false
		}
	}
	return true
}
