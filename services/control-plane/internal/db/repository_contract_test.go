package db

import "testing"

func TestRepositoryTypes(t *testing.T) {
	if (&Repository{}).Pool != nil { t.Fatal("expected nil pool") }
	if (&BillingRepository{}).Pool != nil { t.Fatal("expected nil billing pool") }
	if (&NOCRepository{}).Pool != nil { t.Fatal("expected nil noc pool") }
}
