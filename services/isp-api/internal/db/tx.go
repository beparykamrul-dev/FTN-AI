package db

import "context"

type Tx interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

type TxManager interface {
	Begin(context.Context) (Tx, error)
}

func WithTx(ctx context.Context, manager TxManager, fn func(Tx) error) error {
	tx, err := manager.Begin(ctx)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
