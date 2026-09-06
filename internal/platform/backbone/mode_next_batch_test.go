package backbone

import "testing"

func TestBackboneValidRequiresDistinctEndpoints(t *testing.T) {
	b := Backbone{PrimaryID:"a", SecondaryID:"b", Mode:ModeActive}
	if !b.Valid() { t.Fatal("valid backbone rejected") }
	b.SecondaryID = b.PrimaryID
	if b.Valid() { t.Fatal("duplicate endpoints accepted") }
}
