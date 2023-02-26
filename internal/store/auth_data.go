package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/model"
)

// CreateAuthDataTx creates model.AuthData in transaction
func (s *Storage) CreateAuthDataTx(ctx context.Context, tx Tx, record *model.AuthData) error {
	return s.createAuthData(ctx, tx, record)
}

// createAuthData creates model.AuthData
func (s *Storage) createAuthData(ctx context.Context, db Querier, record *model.AuthData) error {
	query := psql().
		Insert(record.DbTable()).
		SetMap(authDataAttrs(record)).
		Suffix("RETURNING created_at, updated_at")

	err := Getx(ctx, db, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// UpdateAuthDataPassword updates password in model.AuthData
func (s *Storage) UpdateAuthDataPassword(ctx context.Context, record *model.AuthData) error {
	query := psql().
		Update(record.DbTable()).
		Where(sq.Eq{
			"source_id":   record.SourceID,
			"source_type": record.SourceType,
		}).
		Set("password", record.Password).
		Suffix("RETURNING updated_at")

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// GetAuthDataBySourceIDAndType get model.AuthData by source_id, source_type
func (s *Storage) GetAuthDataBySourceIDAndType(
	ctx context.Context,
	sourceID uuid.UUID,
	sourceType model.AuthDataSourceType,
) (*model.AuthData, error) {
	record := &model.AuthData{}
	query := psql().
		Select(asteriskAuthData).
		From(record.DbTable()).
		Where(sq.Eq{
			"source_id":   sourceID,
			"source_type": sourceType,
		}).
		Limit(1)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return nil, dbError(err)
	}

	return record, nil
}

func authDataAttrs(record *model.AuthData) map[string]interface{} {
	return map[string]interface{}{
		"source_id":   record.SourceID,
		"source_type": record.SourceType,
		"password":    record.Password,
	}
}
