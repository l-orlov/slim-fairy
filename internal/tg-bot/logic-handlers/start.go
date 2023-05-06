package logic_handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/pkg/errors"
)

// startMsgTemplate contains start info for chatbot
const startMsgTemplate = `
Привет! Я @%s. Помогу тебе составить план питания.
Ниже список доступных команд.

/start - Начать диалог
/register - Пройти регистрацию
/cancelreg - Прервать регистрацию
/getdietfromai - Получить диету от ИИ`

// Start introduces the bot.
func (h *LogicHandlers) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf(startMsgTemplate, b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send start message")
	}

	return nil
}
