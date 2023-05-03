package logic_handlers

import (
	"context"
	"log"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/pkg/errors"
)

const (
	menuFileName = "menu.txt"
)

// TODO: fix. get user params from db and fill if no mocks
func (h *LogicHandlers) GetDietFromAI(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Printf("user id: %v", ctx.EffectiveSender.User.Id)

	executionCtx := context.Background()
	menuStr, err := h.menuGetter.GetDietByParams(executionCtx, model.GetDietParams{})
	if err != nil {
		return errors.Wrap(err, "h.menuGetter.GetDietByParams")
	}

	reader := strings.NewReader(menuStr)

	_, err = b.SendDocument(ctx.EffectiveChat.Id, gotgbot.NamedFile{
		File:     reader,
		FileName: menuFileName,
	}, &gotgbot.SendDocumentOpts{
		Caption:          "Вот диета для тебя от ИИ (Искусственного Интеллекта)",
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	})
	if err != nil {
		return errors.Wrap(err, "b.SendDocument")
	}

	return nil
}
