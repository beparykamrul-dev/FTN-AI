package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

type PGRepository struct{ db Querier }

func NewPGRepository(db Querier) *PGRepository { return &PGRepository{db: db} }

func (r *PGRepository) Get(ctx context.Context, id string) (Account, error) {
	const q = `SELECT id, name, status FROM accounts WHERE id = $1`
	var a Account
	if err := r.db.QueryRow(ctx, q, id).Scan(&a.ID, &a.Name, &a.Status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return Account{}, fmt.Errorf("account not found: %w", err) }
		return Account{}, err
	}
	return a, nil
}
