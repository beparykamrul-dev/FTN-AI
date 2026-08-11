package module

import "fmt"

type Loader struct { registry *Registry }

func NewLoader(r *Registry) *Loader { return &Loader{registry: r} }

// Load validates and registers a module definition. Runtime construction is
// deliberately separate so modules cannot execute production actions merely
// by being discovered by the Control Plane.
func (l *Loader) Load(m Module) error {
	if m == nil { return fmt.Errorf("module is nil") }
	d := m.Definition()
	if d.Name == "" { return fmt.Errorf("module name is required") }
	if d.Version == "" { return fmt.Errorf("module version is required") }
	if !l.registry.DependenciesReady(d) { return fmt.Errorf("module %q has unavailable dependencies", d.Name) }
	l.registry.Register(m)
	return nil
}
