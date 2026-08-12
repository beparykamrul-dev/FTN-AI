package provider

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrCapabilityDenied = errors.New("provider capability denied")
	ErrCircuitOpen      = errors.New("provider circuit open")
)

type Capability string

const (
	CapabilityDNSRead   Capability = "dns.read"
	CapabilityDNSChange Capability = "dns.change"
)

type Operation struct {
	RequestID       string
	BindingID       string
	ResourceID      string
	Capability      Capability
	Action          string
	DesiredRevision uint64
	IdempotencyKey  string
	Timeout         time.Duration
}

type Result struct {
	RequestID string
	Revision  uint64
	Changed   bool
}

type Adapter interface {
	ID() string
	Capabilities() []Capability
	Health(ctx context.Context) error
	Execute(ctx context.Context, op Operation) (Result, error)
}

type Registry struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
}

func NewRegistry() *Registry { return &Registry{adapters: make(map[string]Adapter)} }

func (r *Registry) Register(a Adapter) error {
	if a == nil || a.ID() == "" { return errors.New("invalid provider adapter") }
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.adapters[a.ID()]; exists { return errors.New("provider adapter already registered") }
	r.adapters[a.ID()] = a
	return nil
}

func (r *Registry) Get(id string) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.adapters[id]
	return a, ok
}

type Policy interface {
	Allow(ctx context.Context, op Operation) bool
}

type Gateway struct {
	Registry *Registry
	Policy   Policy
}

func (g *Gateway) Execute(ctx context.Context, adapterID string, op Operation) (Result, error) {
	if g == nil || g.Registry == nil { return Result{}, errors.New("provider gateway unavailable") }
	if g.Policy != nil && !g.Policy.Allow(ctx, op) { return Result{}, ErrCapabilityDenied }
	a, ok := g.Registry.Get(adapterID)
	if !ok { return Result{}, errors.New("provider adapter not found") }
	allowed := false
	for _, c := range a.Capabilities() { if c == op.Capability { allowed = true; break } }
	if !allowed { return Result{}, ErrCapabilityDenied }
	if err := a.Health(ctx); err != nil { return Result{}, ErrCircuitOpen }
	return a.Execute(ctx, op)
}
