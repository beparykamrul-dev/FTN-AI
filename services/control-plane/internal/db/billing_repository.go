package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BillingRepository struct { Pool *pgxpool.Pool }

func (r *BillingRepository) CountInvoices(ctx context.Context) (int64, error) {
	if r == nil || r.Pool == nil { return 0, nil }
	var count int64
	err := r.Pool.QueryRow(ctx, `SELECT count(*) FROM invoices`).Scan(&count)
	return count, err
}
