package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// TxBeginner is the subset of *pgxpool.Pool needed to start transactions.
type TxBeginner interface {
	Begin(context.Context) (pgx.Tx, error)
}

// WithTx runs fn inside a database transaction: it begins a transaction,
// calls fn with a Queries bound to it, and commits if fn succeeds. If fn
// returns an error or panics, the transaction is rolled back.
func WithTx(ctx context.Context, beginner TxBeginner, queries *Queries, fn func(*Queries) error) error {
	tx, err := beginner.Begin(ctx)
	if err != nil {
		return err
	}
	// Rollback after a successful Commit is a harmless no-op (pgx.ErrTxClosed).
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(queries.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
