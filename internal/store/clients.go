package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/model"
)

// CreateClient creates model.Client
func (s *Storage) CreateClient(ctx context.Context, record *model.Client) error {
	return s.createClient(ctx, s.pool, record)
}

// CreateClientTx creates model.Client in transaction
func (s *Storage) CreateClientTx(ctx context.Context, tx Tx, record *model.Client) error {
	return s.createClient(ctx, tx, record)
}

// createClient creates model.Client
func (s *Storage) createClient(ctx context.Context, db Querier, record *model.Client) error {
	query := psql().
		Insert(record.DbTable()).
		SetMap(clientAttrs(record)).
		Suffix("RETURNING id, created_at, updated_at")

	err := Getx(ctx, db, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// UpdateClient updates model.Client
func (s *Storage) UpdateClient(ctx context.Context, record *model.Client) error {
	query := psql().
		Update(record.DbTable()).
		Where(sq.Eq{"id": record.ID}).
		SetMap(clientAttrs(record)).
		Suffix("RETURNING updated_at")

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// GetClientByID get model.Client by id
func (s *Storage) GetClientByID(ctx context.Context, id uuid.UUID) (*model.Client, error) {
	record := &model.Client{}
	query := psql().
		Select(asteriskClients).
		From(record.DbTable()).
		Where(sq.Eq{"id": id}).
		Limit(1)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return nil, dbError(err)
	}

	return record, nil
}

// GetClientByEmail get model.Client by email
func (s *Storage) GetClientByEmail(ctx context.Context, email string) (*model.Client, error) {
	record := &model.Client{}
	query := psql().
		Select(asteriskClients).
		From(record.DbTable()).
		Where(sq.Eq{"email": email}).
		Limit(1)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return nil, dbError(err)
	}

	return record, nil
}

func clientAttrs(record *model.Client) map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"email":  record.Email,
		"phone":  record.Phone,
		"age":    record.Age,
		"weight": record.Weight,
		"height": record.Height,
		"gender": record.Gender,
	}
}
