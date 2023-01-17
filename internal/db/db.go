package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const pgDSN = "host=127.0.0.1 port=5432 dbname=test sslmode=disable"

type Database struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.New(ctx, pgDSN)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.New")
	}

	return &Database{
		pool: pool,
	}, nil
}

func (db *Database) Close() {
	db.pool.Close()
}
