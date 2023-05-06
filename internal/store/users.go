package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/model"
)

// CreateUser creates model.User
func (s *Storage) CreateUser(ctx context.Context, record *model.User) error {
	return s.createUser(ctx, s.pool, record)
}

// CreateUserTx creates model.User in transaction
func (s *Storage) CreateUserTx(ctx context.Context, tx Tx, record *model.User) error {
	return s.createUser(ctx, tx, record)
}

// createUser creates model.User
func (s *Storage) createUser(ctx context.Context, db Querier, record *model.User) error {
	query := psql().
		Insert(record.DbTable()).
		SetMap(userAttrs(record)).
		Suffix("RETURNING " + asteriskUsers)

	err := Getx(ctx, db, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// UpdateUser updates model.User
func (s *Storage) UpdateUser(ctx context.Context, record *model.User) error {
	query := psql().
		Update(record.DbTable()).
		Where(sq.Eq{"id": record.ID}).
		SetMap(userAttrs(record)).
		Suffix("RETURNING " + asteriskUsers)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// GetUserByID gets model.User by id
func (s *Storage) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	record := &model.User{}
	query := psql().
		Select(asteriskUsers).
		From(record.DbTable()).
		Where(sq.Eq{"id": id}).
		Limit(1)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return nil, dbError(err)
	}

	return record, nil
}

// GetUserByEmail gets model.User by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	record := &model.User{}
	query := psql().
		Select(asteriskUsers).
		From(record.DbTable()).
		Where(sq.Eq{"email": email}).
		Limit(1)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return nil, dbError(err)
	}

	return record, nil
}

func userAttrs(record *model.User) map[string]interface{} {
	return map[string]interface{}{
		"name":        record.Name,
		"email":       record.Email,
		"phone":       record.Phone,
		"telegram_id": record.TelegramID,
		"age":         record.Age,
		"weight":      record.Weight,
		"height":      record.Height,
		"gender":      record.Gender,
		"created_by":  record.CreatedBy,
	}
}
