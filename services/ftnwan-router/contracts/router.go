package router

import "context"

// PacketPlane identifies the implementation used for forwarding. The FTN
// control plane remains independent from the selected dataplane backend.
type PacketPlane string

const (
	PlaneKernel PacketPlane = "kernel"
	PlaneVPP    PacketPlane = "vpp"
	PlaneDPDK   PacketPlane = "dpdk"
)

type Route struct {
	Prefix     string
	NextHop    string
	Interface  string
	Metric     uint32
	Protocol   string
	Installed  bool
}

type Interface struct {
	Name       string
	MAC        string
	MTU        uint32
	AdminUp    bool
	OperUp     bool
	IPv4       []string
	IPv6       []string
}

type RouterState struct {
	NodeID       string
	Plane        PacketPlane
	Interfaces   []Interface
	Routes       []Route
	BGPEnabled   bool
	PPPoEEnabled bool
	NATEnabled   bool
	QoSEnabled   bool
	Conntrack    bool
}

// Dataplane is deliberately narrow: concrete VPP/DPDK/kernel adapters own
// implementation details, while the control plane consumes normalized state.
type Dataplane interface {
	Name() PacketPlane
	Ready(context.Context) bool
	Interfaces(context.Context) ([]Interface, error)
	Routes(context.Context) ([]Route, error)
	ApplyRoute(context.Context, Route) error
	DeleteRoute(context.Context, Route) error
}

// RoutingControl owns validated routing mutations. It does not accept raw
// shell commands or credentials from callers.
type RoutingControl interface {
	PlanRouteChange(context.Context, Route) (string, error)
	ApplyApprovedPlan(context.Context, string) error
	Rollback(context.Context, string) error
}
