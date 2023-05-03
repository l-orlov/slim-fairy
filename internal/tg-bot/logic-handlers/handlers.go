package logic_handlers

import (
	"context"

	"github.com/l-orlov/slim-fairy/internal/model"
)

type (
	DietGetter interface {
		GetDietByParams(ctx context.Context, params model.GetDietParams) (string, error)
	}

	LogicHandlers struct {
		menuGetter DietGetter
	}
)

func New(menuGetter DietGetter) *LogicHandlers {
	return &LogicHandlers{
		menuGetter: menuGetter,
	}
}
