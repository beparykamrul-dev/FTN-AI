package backbone

import "testing"

func TestBackboneValidityRequiresEndpointAndMode(t *testing.T) {
	b := Backbone{PrimaryID:"p", SecondaryID:"s", Mode:ModeActive, Healthy:true}
	if !b.Valid() { t.Fatal("valid backbone rejected") }
	b.SecondaryID = ""
	if b.Valid() { t.Fatal("backbone without secondary endpoint must be invalid") }
}

func TestBackboneCanFailoverRequiresHealthyActiveState(t *testing.T) {
	b := Backbone{PrimaryID:"p", SecondaryID:"s", Mode:ModeActive, Healthy:true}
	if !b.CanFailover() { t.Fatal("healthy active backbone should permit failover proposal") }
	b.Healthy = false
	if b.CanFailover() { t.Fatal("unhealthy backbone must not permit failover") }
}
