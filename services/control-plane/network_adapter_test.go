package main

import (
	"context"
	"testing"
)

type testAdapter struct{}

func (testAdapter) Protocol() string { return " SNMP " }
func (testAdapter) Capabilities(context.Context, NetworkDevice) ([]string, error) {
	return []string{"interfaces"}, nil
}
func (testAdapter) CollectInterfaceState(context.Context, NetworkDevice) ([]InterfaceState, error) {
	return nil, nil
}
func (testAdapter) CollectRoutingState(context.Context, NetworkDevice) ([]RoutingState, error) {
	return nil, nil
}

type testFlowCollector struct{}

func (testFlowCollector) Protocol() string { return "IPFIX" }
func (testFlowCollector) Collect(context.Context, []byte) ([]FlowRecord, error) {
	return nil, nil
}

func TestAdapterRegistryNormalizesProtocols(t *testing.T) {
	r := NewAdapterRegistry()
	r.RegisterAdapter(testAdapter{})
	r.RegisterFlowCollector(testFlowCollector{})
	if _, ok := r.Adapter("snmp"); !ok {
		t.Fatal("SNMP adapter not registered")
	}
	if _, ok := r.FlowCollector("ipfix"); !ok {
		t.Fatal("IPFIX collector not registered")
	}
}
