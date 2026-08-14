package router

import "testing"

func TestPacketPlaneValues(t *testing.T) {
	if PlaneKernel == "" || PlaneVPP == "" || PlaneDPDK == "" {
		t.Fatal("dataplane identifiers must be non-empty")
	}
}

func TestRouterStateDefaultsAreSafe(t *testing.T) {
	state := RouterState{}
	if state.BGPEnabled || state.PPPoEEnabled || state.NATEnabled || state.QoSEnabled || state.Conntrack {
		t.Fatal("router capabilities must not be enabled by zero-value state")
	}
}
