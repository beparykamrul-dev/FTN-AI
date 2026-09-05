package module

import (
	"fmt"
	"strings"
	"sync"
)

type State string
const ( StateRegistered State = "registered"; StateLoaded State = "loaded"; StateStarted State = "started"; StateHealthy State = "healthy"; StateStopped State = "stopped" )
type Lifecycle struct { mu sync.RWMutex; states map[string]State }
func NewLifecycle() *Lifecycle { return &Lifecycle{states: make(map[string]State)} }
func validState(s State) bool { switch s { case StateRegistered, StateLoaded, StateStarted, StateHealthy, StateStopped: return true; default: return false } }
func (l *Lifecycle) Set(name string, state State) error {
	if l == nil { return fmt.Errorf("lifecycle is required") }
	name = strings.TrimSpace(name); state = State(strings.TrimSpace(string(state)))
	if name == "" { return fmt.Errorf("module name is required") }
	if !validState(state) { return fmt.Errorf("invalid module state: %s", state) }
	l.mu.Lock(); defer l.mu.Unlock(); if l.states == nil { l.states = make(map[string]State) }; l.states[name] = state; return nil
}
func (l *Lifecycle) State(name string) (State, bool) {
	if l == nil { return "", false }
	name = strings.TrimSpace(name); if name == "" { return "", false }
	l.mu.RLock(); defer l.mu.RUnlock(); s, ok := l.states[name]; return s, ok
}
