package ftnmesh

import (
	"context"
	"strings"
)

type Node struct { ID string `json:"id"`; Name string `json:"name"`; Role string `json:"role"`; Addresses []string `json:"addresses,omitempty"`; MTU uint32 `json:"mtu,omitempty"`; Ready bool `json:"ready"` }
type Link struct { ID string `json:"id"`; A string `json:"a"`; B string `json:"b"`; CapacityMbps uint64 `json:"capacity_mbps"`; LatencyMs float64 `json:"latency_ms"`; LossPct float64 `json:"loss_pct"`; AdminUp bool `json:"admin_up"`; OperUp bool `json:"oper_up"` }
type RouteIntent struct { ID string `json:"id"`; Prefix string `json:"prefix"`; Source string `json:"source"`; Targets []string `json:"targets"`; Metric uint32 `json:"metric,omitempty"`; Approved bool `json:"approved"` }
type MeshSnapshot struct { Nodes []Node `json:"nodes"`; Links []Link `json:"links"`; Routes []RouteIntent `json:"routes"` }
type MeshAdapter interface { Name() string; Snapshot(context.Context) (MeshSnapshot,error); PlanRoute(context.Context,RouteIntent) (string,error); ApplyApprovedPlan(context.Context,string) error; Rollback(context.Context,string) error }

func (s MeshSnapshot) Validate() error {
	seen := make(map[string]struct{}, len(s.Nodes))
	for _, n := range s.Nodes { id:=strings.TrimSpace(n.ID); if id=="" { return errInvalid("node id is required") }; if _,ok:=seen[id];ok{return errInvalid("duplicate node id: "+id)}; seen[id]=struct{}{} }
	for _, l := range s.Links { id,a,b:=strings.TrimSpace(l.ID),strings.TrimSpace(l.A),strings.TrimSpace(l.B); if id==""||a==""||b==""{return errInvalid("link id and endpoints are required")}; if a==b{return errInvalid("self-link is not allowed: "+id)}; if l.CapacityMbps==0{return errInvalid("link capacity must be positive: "+id)}; if l.LatencyMs<0{return errInvalid("link latency cannot be negative: "+id)}; if l.LossPct<0||l.LossPct>100{return errInvalid("link loss must be between 0 and 100: "+id)} }
	for _, r := range s.Routes { id,prefix,source:=strings.TrimSpace(r.ID),strings.TrimSpace(r.Prefix),strings.TrimSpace(r.Source); if id==""||prefix==""||source==""{return errInvalid("route id, prefix and source are required")}; if r.Approved&&len(r.Targets)==0{return errInvalid("approved route requires a target: "+id)} }
	return nil
}
type validationError string
func(e validationError)Error()string{return string(e)}
func errInvalid(s string)error{return validationError(s)}
