package backbone

import "testing"

func TestBackboneValidityRequiresDistinctEndpoints(t *testing.T) {
	b := Backbone{PrimaryID:"p", SecondaryID:"s", Mode:ModeActive, Healthy:true}
	if !b.Valid() { t.Fatal("valid backbone rejected") }
	b.SecondaryID = b.PrimaryID
	if b.Valid() { t.Fatal("backbone must reject identical endpoints") }
}
