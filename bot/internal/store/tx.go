package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Tx is pgx transaction interface
type Tx pgx.Tx

// WithTransaction executes function in transaction with default isolation level
func (s *Storage) WithTransaction(ctx context.Context, fn func(tx Tx) error) (err error) {
	return s.WithTx(ctx, pgx.TxOptions{}, fn)
}

// WithTx executes function in transaction
func (s *Storage) WithTx(ctx context.Context, options pgx.TxOptions, fn func(tx Tx) error) (err error) {
	tx, err := s.pool.BeginTx(ctx, options)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// panic occurred => rollback and repanic
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			// something went wrong => rollback
			_ = tx.Rollback(ctx)
		} else {
			// ok => commit
			err = tx.Commit(ctx)
		}
	}()

	err = fn(tx)

	return
}
