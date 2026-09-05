package module

import (
	"fmt"
	"strings"
)

type Loader struct{ registry *Registry }
func NewLoader(r *Registry) *Loader { return &Loader{registry: r} }
func (l *Loader) Load(m Module) error {
	if l == nil || l.registry == nil { return fmt.Errorf("module registry is required") }
	if m == nil { return fmt.Errorf("module is nil") }
	d := m.Definition()
	d.Name = strings.TrimSpace(d.Name)
	d.Version = strings.TrimSpace(d.Version)
	if d.Name == "" { return fmt.Errorf("module name is required") }
	if d.Version == "" { return fmt.Errorf("module version is required") }
	if !l.registry.DependenciesReady(d) { return fmt.Errorf("module %q has unavailable dependencies", d.Name) }
	l.registry.Register(m)
	return nil
}
