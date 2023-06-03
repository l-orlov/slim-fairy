package logic_handlers

import (
	"context"

	"github.com/l-orlov/slim-fairy/bot/internal/store"
)

type (
	AIClient interface {
		SendRequest(ctx context.Context, prompt string) (resp string, err error)
	}

	LogicHandlers struct {
		aiClient AIClient
		storage  *store.Storage
	}
)

func New(
	aiClient AIClient,
	storage *store.Storage,
) *LogicHandlers {
	return &LogicHandlers{
		aiClient: aiClient,
		storage:  storage,
	}
}
