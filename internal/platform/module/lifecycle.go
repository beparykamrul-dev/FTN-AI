package module

import "fmt"

type State string

const (
	StateRegistered State = "registered"
	StateLoaded     State = "loaded"
	StateStarted    State = "started"
	StateHealthy    State = "healthy"
	StateStopped    State = "stopped"
)

type Lifecycle struct {
	states map[string]State
}

func NewLifecycle() *Lifecycle { return &Lifecycle{states: make(map[string]State)} }

func (l *Lifecycle) Set(name string, state State) error {
	if name == "" {
		return fmt.Errorf("module name is required")
	}
	if state == "" {
		return fmt.Errorf("module state is required")
	}
	l.states[name] = state
	return nil
}

func (l *Lifecycle) State(name string) (State, bool) {
	s, ok := l.states[name]
	return s, ok
}
