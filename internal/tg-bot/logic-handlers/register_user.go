package logic_handlers

import (
	"fmt"
	"html"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/pkg/errors"
)

// TODO: add registration handlers

// These are the user registration handler states.
const (
	RegisterName             = "name"
	RegisterAge              = "age"
	RegisterWeight           = "weight"
	RegisterHeight           = "height"
	RegisterGender           = "gender"
	RegisterPhysicalActivity = "physicalactivity"
)

// StartUserRegistration starts user registration
func (h *LogicHandlers) StartUserRegistration(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s.\nWhat is your name?.", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send start message")
	}
	return handlers.NextConversationState(RegisterName)
}

// CancelUserRegistration cancels user registration
func (h *LogicHandlers) CancelUserRegistration(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Oh, goodbye!", &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send cancel message")
	}
	return handlers.EndConversation()
}

// name gets the user's name
func (h *LogicHandlers) RegisterUserName(b *gotgbot.Bot, ctx *ext.Context) error {
	inputName := ctx.EffectiveMessage.Text
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Nice to meet you, %s!\n\nAnd how old are you?", html.EscapeString(inputName)), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send name message")
	}
	return handlers.NextConversationState(RegisterAge)
}

// age gets the user's age
func (h *LogicHandlers) RegisterUserAge(b *gotgbot.Bot, ctx *ext.Context) error {
	inputAge := ctx.EffectiveMessage.Text
	ageNumber, err := strconv.ParseInt(inputAge, 10, 64)
	if err != nil {
		// If the number is not valid, try again!
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("This doesn't seem to be a number. Could you repeat?"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		// We try the age handler again
		return handlers.NextConversationState(RegisterAge)
	}

	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Ah, you're %d years old!", ageNumber), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send age message")
	}
	return handlers.EndConversation()
}
