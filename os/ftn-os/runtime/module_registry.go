package runtime

import (
	"fmt"
	"sync"
)

type Module struct {
	Name         string
	Version      string
	Dependencies []string
}

type ModuleRegistry struct {
	mu      sync.RWMutex
	modules map[string]Module
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{modules: make(map[string]Module)}
}

func (r *ModuleRegistry) Register(module Module) error {
	if module.Name == "" || module.Version == "" {
		return fmt.Errorf("module name and version are required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.modules[module.Name]; exists {
		return fmt.Errorf("module already registered: %s", module.Name)
	}
	r.modules[module.Name] = module
	return nil
}

func (r *ModuleRegistry) Get(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	module, ok := r.modules[name]
	return module, ok
}

func (r *ModuleRegistry) List() []Module {
	r.mu.RLock()
	defer r.mu.RUnlock()
	modules := make([]Module, 0, len(r.modules))
	for _, module := range r.modules {
		modules = append(modules, module)
	}
	return modules
}
