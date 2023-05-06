package store

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/l-orlov/slim-fairy/internal/config"
	"github.com/pkg/errors"
)

const (
	//defaultTimeout = 3 * time.Second
	defaultTimeout = 30 * time.Minute // TODO: fix
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context) (*Storage, error) {
	pgDSN := config.Get().PgDSN
	pool, err := pgxpool.New(ctx, pgDSN)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.New")
	}

	return &Storage{
		pool: pool,
	}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

// Postgres specific squirrel builder
func psql() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

// Select executes select query
func Select(ctx context.Context, db pgxscan.Querier, dst interface{}, query string, args ...interface{}) error {
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return pgxscan.Select(dbCtx, db, dst, query, args...)
}

// Selectx executes select query with squirrel.Sqlizer
func Selectx(ctx context.Context, db pgxscan.Querier, dst interface{}, sqlizer sq.Sqlizer) error {
	stmt, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}

	return Select(ctx, db, dst, stmt, args...)
}

// Querier can execute sql query and get the pgx.Rows
type Querier interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

// Get executes get query
func Get(ctx context.Context, db Querier, dst interface{}, query string, args ...interface{}) error {
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return pgxscan.Get(dbCtx, db, dst, query, args...)
}

// Getx executes get query with squirrel.Sqlizer
func Getx(ctx context.Context, db Querier, dst interface{}, sqlizer sq.Sqlizer) error {
	stmt, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}

	return Get(ctx, db, dst, stmt, args...)
}

// Executor can execute sql query
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// Exec executes query
func Exec(ctx context.Context, db Executor, query string, args ...interface{}) (pgconn.CommandTag, error) {
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return db.Exec(dbCtx, query, args...)
}

// Execx executes query with squirrel.Sqlizer
func Execx(ctx context.Context, db Executor, sqlizer sq.Sqlizer) (pgconn.CommandTag, error) {
	stmt, args, err := sqlizer.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return Exec(ctx, db, stmt, args...)
}
