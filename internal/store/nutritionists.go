package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/model"
)

// CreateNutritionist creates model.Nutritionist
func (s *Storage) CreateNutritionist(ctx context.Context, record *model.Nutritionist) error {
	query := psql().
		Insert(record.DbTable()).
		SetMap(nutritionistAttrs(record)).
		Suffix("RETURNING id, created_at, updated_at")

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// UpdateNutritionist updates model.Nutritionist
func (s *Storage) UpdateNutritionist(ctx context.Context, record *model.Nutritionist) error {
	query := psql().
		Update(record.DbTable()).
		Where(sq.Eq{"id": record.ID}).
		SetMap(nutritionistAttrs(record)).
		Suffix("RETURNING updated_at")

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// GetNutritionistByID get model.Nutritionist by id
func (s *Storage) GetNutritionistByID(ctx context.Context, id uuid.UUID) (*model.Nutritionist, error) {
	record := &model.Nutritionist{}
	query := psql().
		Select(asteriskNutritionists).
		From(record.DbTable()).
		Where(sq.Eq{"id": record.ID}).
		Limit(1)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return nil, dbError(err)
	}

	return record, nil
}

func nutritionistAttrs(record *model.Nutritionist) map[string]interface{} {
	return map[string]interface{}{
		"name":   record.Name,
		"email":  record.Email,
		"phone":  record.Phone,
		"age":    record.Age,
		"gender": record.Gender,
		"info":   record.Info,
	}
}
