package db

import (
	"context"
	"testing"
)

func TestBillingRepositoryNilPool(t *testing.T) {
	got, err := (&BillingRepository{}).CountInvoices(context.Background())
	if err != nil || got != 0 { t.Fatalf("got %d, err %v", got, err) }
}

func TestNOCRepositoryNilPool(t *testing.T) {
	got, err := (&NOCRepository{}).CountNodes(context.Background())
	if err != nil || got != 0 { t.Fatalf("got %d, err %v", got, err) }
}
