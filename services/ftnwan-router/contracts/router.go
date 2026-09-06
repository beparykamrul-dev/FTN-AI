package router

import "context"

type PacketPlane string
const ( PlaneKernel PacketPlane = "kernel"; PlaneVPP PacketPlane = "vpp"; PlaneDPDK PacketPlane = "dpdk" )
type Route struct { Prefix string; NextHop string; Interface string; Metric uint32; Protocol string; Installed bool }
type Interface struct { Name string; MAC string; MTU uint32; AdminUp bool; OperUp bool; IPv4 []string; IPv6 []string }
type RouterState struct { NodeID string; Plane PacketPlane; Interfaces []Interface; Routes []Route; BGPEnabled bool; PPPoEEnabled bool; NATEnabled bool; QoSEnabled bool; Conntrack bool }
type Dataplane interface { Name() PacketPlane; Ready(context.Context) bool; Interfaces(context.Context)([]Interface,error); Routes(context.Context)([]Route,error); ApplyRoute(context.Context,Route)error; DeleteRoute(context.Context,Route)error }
type RoutingControl interface { PlanRouteChange(context.Context,Route)(string,error); ApplyApprovedPlan(context.Context,string)error; Rollback(context.Context,string)error }
