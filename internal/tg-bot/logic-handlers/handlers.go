package logic_handlers

import (
	"context"

	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
)

type (
	DietGetter interface {
		GetDietByParams(ctx context.Context, params model.GetDietParams) (string, error)
	}

	LogicHandlers struct {
		menuGetter DietGetter
		storage    *store.Storage
	}
)

func New(
	menuGetter DietGetter,
	storage *store.Storage,
) *LogicHandlers {
	return &LogicHandlers{
		menuGetter: menuGetter,
		storage:    storage,
	}
}
