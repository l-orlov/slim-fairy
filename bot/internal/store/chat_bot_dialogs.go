package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/bot/internal/model"
)

// CreateChatBotDialog creates model.ChatBotDialog
func (s *Storage) CreateChatBotDialog(ctx context.Context, record *model.ChatBotDialog) error {
	query := psql().
		Insert(record.DbTable()).
		SetMap(chatBotDialogAttrs(record)).
		Suffix("RETURNING " + asteriskChatBotDialogs)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// UpdateChatBotDialog updates model.ChatBotDialog
func (s *Storage) UpdateChatBotDialog(ctx context.Context, record *model.ChatBotDialog) error {
	query := psql().
		Update(record.DbTable()).
		Where(sq.Eq{"id": record.ID}).
		SetMap(chatBotDialogAttrs(record)).
		Suffix("RETURNING " + asteriskChatBotDialogs)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// UpdateChatBotDialogStatusTx updates status for model.ChatBotDialog in transaction
func (s *Storage) UpdateChatBotDialogStatusTx(ctx context.Context, tx Tx, status model.ChatBotDialogStatus, id uuid.UUID) error {
	query := psql().
		Update(model.ChatBotDialog{}.DbTable()).
		Where(sq.Eq{"id": id}).
		Set("status", status)

	_, err := Execx(ctx, tx, query)
	if err != nil {
		return dbError(err)
	}
	return nil
}

// GetChatBotDialogByKeyFields gets last model.ChatBotDialog by user_telegram_id, kind, status
func (s *Storage) GetChatBotDialogByKeyFields(
	ctx context.Context, userTelegramID int64,
	kind model.ChatBotDialogKind, status model.ChatBotDialogStatus,
) (*model.ChatBotDialog, error) {
	record := &model.ChatBotDialog{}
	query := psql().
		Select(asteriskChatBotDialogs).
		From(record.DbTable()).
		Where(sq.Eq{
			"user_telegram_id": userTelegramID,
			"kind":             kind,
			"status":           status,
		}).
		OrderBy("created_at DESC").
		Limit(1)

	err := Getx(ctx, s.pool, record, query)
	if err != nil {
		return nil, dbError(err)
	}

	return record, nil
}

func chatBotDialogAttrs(record *model.ChatBotDialog) map[string]interface{} {
	return map[string]interface{}{
		"user_telegram_id": record.UserTelegramID,
		"kind":             record.Kind,
		"status":           record.Status,
		"data":             record.DataJSON,
	}
}
