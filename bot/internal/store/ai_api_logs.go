package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/l-orlov/slim-fairy/bot/internal/model"
)

// CreateAIAPILog creates model.AIAPILogs
func (s *Storage) CreateAIAPILog(ctx context.Context, record *model.AIAPILog) error {
	query := psql().
		Insert(record.DbTable()).
		SetMap(aiAPILogAttrs(record)).
		Suffix("RETURNING " + asteriskAIAPILogs)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// UpdateAIAPILog updates model.AIAPILog
func (s *Storage) UpdateAIAPILog(ctx context.Context, record *model.AIAPILog) error {
	query := psql().
		Update(record.DbTable()).
		Where(sq.Eq{"id": record.ID}).
		SetMap(aiAPILogAttrs(record)).
		Suffix("RETURNING " + asteriskChatBotDialogs)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

func aiAPILogAttrs(record *model.AIAPILog) map[string]interface{} {
	return map[string]interface{}{
		"prompt":      record.Prompt,
		"response":    record.Response,
		"user_id":     record.UserID,
		"source_id":   record.SourceID,
		"source_type": record.SourceType,
	}
}
