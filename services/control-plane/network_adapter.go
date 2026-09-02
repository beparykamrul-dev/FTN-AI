package main

import "context"

type NetworkAdapter interface {
	Protocol() string
	Capabilities(context.Context, NetworkDevice) ([]string, error)
	CollectInterfaceState(context.Context, NetworkDevice) ([]InterfaceState, error)
	CollectRoutingState(context.Context, NetworkDevice) ([]RoutingState, error)
}

type FlowCollector interface {
	Protocol() string
	Collect(context.Context, []byte) ([]FlowRecord, error)
}

type AdapterRegistry struct {
	adapters map[string]NetworkAdapter
	flows map[string]FlowCollector
}

func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{adapters: make(map[string]NetworkAdapter), flows: make(map[string]FlowCollector)}
}

func (r *AdapterRegistry) RegisterAdapter(a NetworkAdapter) {
	if a == nil { return }
	r.adapters[normalizeProtocol(a.Protocol())] = a
}

func (r *AdapterRegistry) RegisterFlowCollector(c FlowCollector) {
	if c == nil { return }
	r.flows[normalizeProtocol(c.Protocol())] = c
}

func (r *AdapterRegistry) Adapter(protocol string) (NetworkAdapter, bool) {
	a, ok := r.adapters[normalizeProtocol(protocol)]
	return a, ok
}

func (r *AdapterRegistry) FlowCollector(protocol string) (FlowCollector, bool) {
	c, ok := r.flows[normalizeProtocol(protocol)]
	return c, ok
}
