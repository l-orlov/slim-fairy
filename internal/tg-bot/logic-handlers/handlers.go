package logic_handlers

import (
	"context"

	"github.com/l-orlov/slim-fairy/internal/model"
)

type (
	MenuGetter interface {
		GetMenuByParams(ctx context.Context, params model.GetMenuParams) (string, error)
	}

	LogicHandlers struct {
		menuGetter MenuGetter
	}
)

func New(menuGetter MenuGetter) *LogicHandlers {
	return &LogicHandlers{
		menuGetter: menuGetter,
	}
}
