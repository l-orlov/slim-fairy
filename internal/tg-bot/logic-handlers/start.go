package logic_handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// startMsgTemplate contains start info for chatbot
const startMsgTemplate = `
Привет! Я @%s. Помогу вам составить план питания.
Ниже список доступных команд.

/start - Начать диалог
/register - Пройти регистрацию
/cancelreg - Прервать регистрацию
/getdietfromai - Получить диету от ИИ`

// Start introduces the bot.
func (h *LogicHandlers) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	replyMsg := fmt.Sprintf(startMsgTemplate, b.User.Username)
	reply(b, ctx, replyMsg)

	return nil
}
