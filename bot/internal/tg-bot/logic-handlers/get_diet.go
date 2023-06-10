package logic_handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// StartGettingDiet .
func (h *LogicHandlers) StartGettingDiet(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := `На данный момент функционал находится в разработке`
	sendMessage(b, ctx, msg, nil)

	return nil
}
