package device

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu      sync.RWMutex
	adapters map[Kind]Adapter
}

func NewRegistry() *Registry {
	return &Registry{adapters: make(map[Kind]Adapter)}
}

func (r *Registry) Register(adapter Adapter) error {
	if adapter == nil || adapter.Kind() == "" {
		return fmt.Errorf("invalid device adapter")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.adapters[adapter.Kind()]; exists {
		return fmt.Errorf("adapter already registered: %s", adapter.Kind())
	}
	r.adapters[adapter.Kind()] = adapter
	return nil
}

func (r *Registry) Get(kind Kind) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.adapters[kind]
	return a, ok
}
