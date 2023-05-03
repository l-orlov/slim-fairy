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

// TODO: сохранять промежуточные данные в каком-то кэше. после подтверждения записывать в бд

const (
	// User registration handler states
	RegisterName             = "name"
	RegisterAge              = "age"
	RegisterWeight           = "weight"
	RegisterHeight           = "height"
	RegisterGender           = "gender"
	RegisterPhysicalActivity = "physicalactivity"
	RegisterConfirm          = "сonfirm"

	// userRegistrationStartInfo contains start info for user registration
	userRegistrationStartInfo = `
Чтобы точнее составить диету, нужны параметры (рост, вес и другие).
Я буду спрашивать по очереди каждый параметр.
Отвечай мне в следующем сообщении.

Если решишь прервать регистрацию, то используй команду:
/cancelreg

Начнем. Как тебя зовут?`
)

// StartUserRegistration starts user registration
func (h *LogicHandlers) StartUserRegistration(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, userRegistrationStartInfo, &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send start message")
	}
	return handlers.NextConversationState(RegisterName)
}

// CancelUserRegistration cancels user registration
func (h *LogicHandlers) CancelUserRegistration(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "ОК, тогда в следующий раз", &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send cancel message")
	}
	return handlers.EndConversation()
}

// RegisterUserName registers user name
func (h *LogicHandlers) RegisterUserName(b *gotgbot.Bot, ctx *ext.Context) error {
	input := ctx.EffectiveMessage.Text

	// Validate name
	const validLen = 100
	if len(input) > validLen {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Слишком длинное имя. Сократи до %d символов", validLen), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterName)
	}

	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Рад знакомству, %s!\nСколько тебе лет?", html.EscapeString(input)), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send name message")
	}
	return handlers.NextConversationState(RegisterAge)
}

// RegisterUserAge registers user age
func (h *LogicHandlers) RegisterUserAge(b *gotgbot.Bot, ctx *ext.Context) error {
	input := ctx.EffectiveMessage.Text
	ageNumber, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Нужно число. Попробуй еще раз"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterAge)
	}
	if ageNumber < 0 || ageNumber > 150 {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Число не подходит для возраста. Попробуй еще раз"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterAge)
	}

	_, err = ctx.EffectiveMessage.Reply(b, "Прекрасный возраст.\nСколько ты весишь кг?", &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send age message")
	}
	return handlers.NextConversationState(RegisterWeight)
}

// RegisterUserWeight registers user weight
func (h *LogicHandlers) RegisterUserWeight(b *gotgbot.Bot, ctx *ext.Context) error {
	input := ctx.EffectiveMessage.Text
	weight, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Нужно число. Попробуй еще раз"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterWeight)
	}
	if weight < 0 || weight > 300 {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Число не подходит для веса. Попробуй еще раз"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterWeight)
	}

	_, err = ctx.EffectiveMessage.Reply(b, "Хорошо.\nКакой у тебя рост см?", &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send age message")
	}
	return handlers.NextConversationState(RegisterHeight)
}

// RegisterUserHeight registers user height
func (h *LogicHandlers) RegisterUserHeight(b *gotgbot.Bot, ctx *ext.Context) error {
	input := ctx.EffectiveMessage.Text
	height, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Нужно число. Попробуй еще раз"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterHeight)
	}
	if height < 0 || height > 250 {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Число не подходит для роста. Попробуй еще раз"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterHeight)
	}

	_, err = ctx.EffectiveMessage.Reply(b, "ОК.\nКакого ты пола?\nНапиши: м (мужчина) или ж (женщина)", &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return errors.Wrap(err, "failed to send age message")
	}
	return handlers.NextConversationState(RegisterGender)
}

// RegisterUserGender registers user Gender
func (h *LogicHandlers) RegisterUserGender(b *gotgbot.Bot, ctx *ext.Context) error {
	input := ctx.EffectiveMessage.Text
	if input != "м" && input != "ж" {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Попробуй еще раз. Напиши: м или ж"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterGender)
	}

	_, err := ctx.EffectiveMessage.Reply(b, `
Принято.
Теперь какая у тебя уровень физической активности?
Напиши: н (низкий) или с (средний) или в (высокий)`,
		&gotgbot.SendMessageOpts{
			ParseMode: "html",
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to send age message")
	}
	return handlers.NextConversationState(RegisterPhysicalActivity)
}

// RegisterUserPhysicalActivity registers user physical activity
func (h *LogicHandlers) RegisterUserPhysicalActivity(b *gotgbot.Bot, ctx *ext.Context) error {
	input := ctx.EffectiveMessage.Text
	if input != "н" && input != "с" && input != "в" {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Попробуй еще раз. Напиши: н или с или в"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterPhysicalActivity)
	}

	_, err := ctx.EffectiveMessage.Reply(b, `
Принято.
Тут инфа пользака.
Подтверди, что данные верны.
Напиши: да (верны) или нет (не верны)`,
		&gotgbot.SendMessageOpts{
			ParseMode: "html",
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to send age message")
	}
	return handlers.NextConversationState(RegisterConfirm)
}

// ConfirmUserRegistration confirms user registration
func (h *LogicHandlers) ConfirmUserRegistration(b *gotgbot.Bot, ctx *ext.Context) error {
	input := ctx.EffectiveMessage.Text
	if input != "да" && input != "нет" {
		// Not valid -> try again
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Попробуй еще раз. Напиши: да или нет"), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return handlers.NextConversationState(RegisterConfirm)
	}

	if input == "нет" {
		_, err := ctx.EffectiveMessage.Reply(b, `
Попробуй пройти регистрацию снова:
/register`,
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
			},
		)
		if err != nil {
			return errors.Wrap(err, "failed to send age message")
		}
		return handlers.EndConversation()
	}

	_, err := ctx.EffectiveMessage.Reply(b, `
Ты успешно прошел регистрацию!
Уже учел твои параметры для составления диеты.`,
		&gotgbot.SendMessageOpts{
			ParseMode: "html",
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to send age message")
	}
	return handlers.EndConversation()
}
