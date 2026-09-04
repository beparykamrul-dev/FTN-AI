package ftnmesh

import (
	"context"
	"net"
)

// Node is a normalized backbone/POP mesh participant. Concrete transport and
// routing implementations stay behind the MeshAdapter boundary.
type Node struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Role      string   `json:"role"`
	Addresses []string `json:"addresses,omitempty"`
	MTU       uint32   `json:"mtu,omitempty"`
	Ready     bool     `json:"ready"`
}

type Link struct {
	ID           string  `json:"id"`
	A            string  `json:"a"`
	B            string  `json:"b"`
	CapacityMbps uint64  `json:"capacity_mbps"`
	LatencyMs    float64 `json:"latency_ms"`
	LossPct      float64 `json:"loss_pct"`
	AdminUp      bool    `json:"admin_up"`
	OperUp       bool    `json:"oper_up"`
}

type RouteIntent struct {
	ID       string   `json:"id"`
	Prefix   string   `json:"prefix"`
	Source   string   `json:"source"`
	Targets  []string `json:"targets"`
	Metric   uint32   `json:"metric,omitempty"`
	Approved bool     `json:"approved"`
}

type MeshSnapshot struct {
	Nodes  []Node        `json:"nodes"`
	Links  []Link        `json:"links"`
	Routes []RouteIntent `json:"routes"`
}

// MeshAdapter exposes only normalized state and validated, approval-gated
// mutations. It never accepts arbitrary shell commands or credentials.
type MeshAdapter interface {
	Name() string
	Snapshot(context.Context) (MeshSnapshot, error)
	PlanRoute(context.Context, RouteIntent) (string, error)
	ApplyApprovedPlan(context.Context, string) error
	Rollback(context.Context, string) error
}

// Validate performs deterministic boundary checks before a plan can enter
// the approval/execution pipeline.
func (s MeshSnapshot) Validate() error {
	seen := make(map[string]struct{}, len(s.Nodes))
	for _, n := range s.Nodes {
		if n.ID == "" {
			return errInvalid("node id is required")
		}
		if _, ok := seen[n.ID]; ok {
			return errInvalid("duplicate node id: " + n.ID)
		}
		seen[n.ID] = struct{}{}
	}

	for _, l := range s.Links {
		if l.ID == "" || l.A == "" || l.B == "" {
			return errInvalid("link id and endpoints are required")
		}
		if l.A == l.B {
			return errInvalid("self-link is not allowed: " + l.ID)
		}
		if _, ok := seen[l.A]; !ok {
			return errInvalid("link endpoint does not exist: " + l.A)
		}
		if _, ok := seen[l.B]; !ok {
			return errInvalid("link endpoint does not exist: " + l.B)
		}
		if l.CapacityMbps == 0 {
			return errInvalid("link capacity must be positive: " + l.ID)
		}
		if l.LatencyMs < 0 {
			return errInvalid("link latency cannot be negative: " + l.ID)
		}
		if l.LossPct < 0 || l.LossPct > 100 {
			return errInvalid("link loss must be between 0 and 100: " + l.ID)
		}
	}

	for _, r := range s.Routes {
		if r.ID == "" || r.Prefix == "" || r.Source == "" {
			return errInvalid("route id, prefix and source are required")
		}
		if _, _, err := net.ParseCIDR(r.Prefix); err != nil {
			return errInvalid("invalid route prefix: " + r.Prefix)
		}
		if _, ok := seen[r.Source]; !ok {
			return errInvalid("route source does not exist: " + r.Source)
		}
		for _, target := range r.Targets {
			if target == "" {
				return errInvalid("route target cannot be empty: " + r.ID)
			}
			if _, ok := seen[target]; !ok {
				return errInvalid("route target does not exist: " + target)
			}
			if target == r.Source {
				return errInvalid("route source cannot target itself: " + r.ID)
			}
		}
		if !r.Approved {
			continue
		}
		if len(r.Targets) == 0 {
			return errInvalid("approved route requires a target: " + r.ID)
		}
	}
	return nil
}

type validationError string

func (e validationError) Error() string { return string(e) }

func errInvalid(s string) error { return validationError(s) }
