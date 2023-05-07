package logic_handlers

import (
	"context"

	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
)

type (
	DietGetter interface {
		GetDietByParams(ctx context.Context, params *model.GetDietParams) (string, error)
	}

	LogicHandlers struct {
		dietGetter DietGetter
		storage    *store.Storage
	}
)

func New(
	dietGetter DietGetter,
	storage *store.Storage,
) *LogicHandlers {
	return &LogicHandlers{
		dietGetter: dietGetter,
		storage:    storage,
	}
}
