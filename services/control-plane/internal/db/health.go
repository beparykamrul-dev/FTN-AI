package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Ping(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil { return nil }
	return pool.Ping(ctx)
}
