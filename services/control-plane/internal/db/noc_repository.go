package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NOCRepository struct { Pool *pgxpool.Pool }

func (r *NOCRepository) CountNodes(ctx context.Context) (int64, error) {
	if r == nil || r.Pool == nil { return 0, nil }
	var count int64
	err := r.Pool.QueryRow(ctx, `SELECT count(*) FROM control_nodes`).Scan(&count)
	return count, err
}
