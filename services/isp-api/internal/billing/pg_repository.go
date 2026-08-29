package billing

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

type PGRepository struct{ db Querier }

func NewPGRepository(db Querier) *PGRepository { return &PGRepository{db: db} }

func (r *PGRepository) GetInvoice(ctx context.Context, id string) (Invoice, error) {
	const q = `SELECT id, account_id, amount_minor, currency, status FROM invoices WHERE id = $1`
	var i Invoice
	if err := r.db.QueryRow(ctx, q, id).Scan(&i.ID, &i.AccountID, &i.AmountMinor, &i.Currency, &i.Status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return Invoice{}, fmt.Errorf("invoice not found: %w", err) }
		return Invoice{}, err
	}
	return i, nil
}

func (r *PGRepository) ListAccountInvoices(ctx context.Context, accountID string) ([]Invoice, error) {
	const q = `SELECT id, account_id, amount_minor, currency, status FROM invoices WHERE account_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, accountID)
	if err != nil { return nil, err }
	defer rows.Close()
	out := make([]Invoice, 0)
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(&i.ID, &i.AccountID, &i.AmountMinor, &i.Currency, &i.Status); err != nil { return nil, err }
		out = append(out, i)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return out, nil
}
