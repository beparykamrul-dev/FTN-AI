package router

import (
	"fmt"
	"net"
	"strings"
)

// ValidateRoute performs deterministic, side-effect-free validation before a
// route can enter the FTN routing-control approval pipeline.
func ValidateRoute(r Route) error {
	if strings.TrimSpace(r.Prefix) == "" {
		return fmt.Errorf("route prefix is required")
	}
	if _, _, err := net.ParseCIDR(r.Prefix); err != nil {
		return fmt.Errorf("invalid route prefix %q: %w", r.Prefix, err)
	}
	if strings.TrimSpace(r.NextHop) != "" && net.ParseIP(r.NextHop) == nil {
		return fmt.Errorf("invalid next hop %q", r.NextHop)
	}
	if strings.TrimSpace(r.Interface) == "" && strings.TrimSpace(r.NextHop) == "" {
		return fmt.Errorf("route requires an interface or next hop")
	}
	return nil
}

// ValidateState checks the normalized router state without mutating the
// selected dataplane. It is intended for discovery, reconciliation and
// approval workflows.
func ValidateState(s RouterState) error {
	if strings.TrimSpace(s.NodeID) == "" {
		return fmt.Errorf("node id is required")
	}
	if s.Plane != PlaneKernel && s.Plane != PlaneVPP && s.Plane != PlaneDPDK {
		return fmt.Errorf("unsupported packet plane %q", s.Plane)
	}
	for _, r := range s.Routes {
		if err := ValidateRoute(r); err != nil {
			return err
		}
	}
	return nil
}
