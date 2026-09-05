package agent

import "sync"

type Adapter interface {
	Kind() DeviceKind
	Capabilities() []string
}

type AdapterRegistry struct {
	mu      sync.RWMutex
	adapters map[DeviceKind]Adapter
}

func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{adapters: make(map[DeviceKind]Adapter)}
}

func (r *AdapterRegistry) Register(a Adapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[a.Kind()] = a
}

func (r *AdapterRegistry) Get(kind DeviceKind) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.adapters[kind]
	return a, ok
}

func (r *AdapterRegistry) Kinds() []DeviceKind {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]DeviceKind, 0, len(r.adapters))
	for k := range r.adapters {
		out = append(out, k)
	}
	return out
}
