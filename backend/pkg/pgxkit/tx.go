package pgxkit

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTx runs fn inside a database transaction.
// Commits on nil error; rolls back on any error or panic.
// Context-aware — cancels the tx if ctx is done.
//
// Example:
//
//	err := pgxkit.WithTx(ctx, pool, func(tx pgx.Tx) error {
//	    // ... use tx via sqlc's WithTx(tx) ...
//	})
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("pgxkit: begin tx: %w", err)
	}

	// Rollback is a no-op if tx already committed.
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
