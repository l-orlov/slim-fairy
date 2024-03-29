package logic_handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/l-orlov/slim-fairy/bot/internal/model"
	"github.com/l-orlov/slim-fairy/bot/internal/store"
	"github.com/l-orlov/slim-fairy/bot/pkg/ctxutil"
	"github.com/l-orlov/slim-fairy/bot/pkg/ptrconv"
	"github.com/pkg/errors"
)

// Handler states for getting diet from AI
const (
	GetDietFromAISelectMeals         = "select_meals"
	GetDietFromAISelectSnacks        = "select_snacks"
	GetDietFromAISelectOrderProducts = "select_order_products"
)

// Callback keys for getting diet from AI
const (
	GetDietFromAICbSelectMeals2         = "select_meals_2"
	GetDietFromAICbSelectMeals3         = "select_meals_3"
	GetDietFromAICbSelectSnacks0        = "select_snacks_0"
	GetDietFromAICbSelectSnacks1        = "select_snacks_1"
	GetDietFromAICbSelectSnacks2        = "select_snacks_2"
	GetDietFromAISelectOrderProductsYes = "select_order_products_y"
	GetDietFromAISelectOrderProductsNo  = "select_order_products_n"
)

// Prompts templates
const (
	// Get usual diet
	promptTemplateGetDiet = `
Составь диету на один день с тремя приемами пищи в день и %s.
Укажи список ингредиентов и их количество, КБЖУ для каждого блюда и для всей диеты.

Возраст: %d.
Рост: %d см.
Вес: %d кг.
Пол: %s.
Уровень физической активности: %s.

Составь список продуктов и их количество для покупки.`
	// Get diet for interval fasting
	promptTemplateGetIntervalDiet = `
Составь диету на один день для интервального голодания с двумя приемами пищи в день и без перекусов.
Укажи список ингредиентов и их количество, КБЖУ для каждого блюда и для всей диеты.

Возраст: %d.
Рост: %d см.
Вес: %d кг.
Пол: %s.
Уровень физической активности: %s.

Составь список продуктов и их количество для покупки.`
)

const (
	menuFileName = "menu.txt"
)

// Timeouts
const (
	sendRequestToAITimeout = 2 * time.Minute
	createLogTimeout       = 1 * time.Minute
)

const sbermarketURL = "https://sbermarket.ru/"

/* TODO:
лимит - 7 запросов в 1 день.
А дальше идите в платную версию
*/

// StartGettingDietFromAI .
func (h *LogicHandlers) StartGettingDietFromAI(b *gotgbot.Bot, ctx *ext.Context) error {
	opts := &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "2", CallbackData: GetDietFromAICbSelectMeals2},
				{Text: "3", CallbackData: GetDietFromAICbSelectMeals3},
			}},
		},
	}
	return replyInConversation(b, ctx, "Выбери количество приемов пищи в день", GetDietFromAISelectMeals, opts)
}

// CancelGettingDiet .
func (h *LogicHandlers) CancelGettingDiet(b *gotgbot.Bot, ctx *ext.Context) error {
	return endConversation(b, ctx, "ОК, тогда в следующий раз")
}

// SelectMeals2 .
func (h *LogicHandlers) SelectMeals2(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали 2",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали 2 приема пищи в день", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	return h.sendDietFromAI(b, ctx, model.MealsAndSnacksNumber{
		MealsNumberPerDay:  2,
		SnacksNumberPerDay: 0,
	})
}

// SelectMeals3 .
func (h *LogicHandlers) SelectMeals3(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали 3",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали 3 приема пищи в день", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	opts := &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "0", CallbackData: GetDietFromAICbSelectSnacks0},
				{Text: "1", CallbackData: GetDietFromAICbSelectSnacks1},
				{Text: "2", CallbackData: GetDietFromAICbSelectSnacks2},
			}},
		},
	}
	return replyInConversation(b, ctx, "Выбери количество перекусов в день", GetDietFromAISelectSnacks, opts)
}

// SelectSnacks0 .
func (h *LogicHandlers) SelectSnacks0(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали 0",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали 0 перекусов в день", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	return h.sendDietFromAI(b, ctx, model.MealsAndSnacksNumber{
		MealsNumberPerDay:  3,
		SnacksNumberPerDay: 0,
	})
}

// SelectSnacks1 .
func (h *LogicHandlers) SelectSnacks1(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали 1",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали 1 перекус в день", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	return h.sendDietFromAI(b, ctx, model.MealsAndSnacksNumber{
		MealsNumberPerDay:  3,
		SnacksNumberPerDay: 1,
	})
}

// SelectSnacks2 .
func (h *LogicHandlers) SelectSnacks2(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали 2",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали 2 перекуса в день", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	return h.sendDietFromAI(b, ctx, model.MealsAndSnacksNumber{
		MealsNumberPerDay:  3,
		SnacksNumberPerDay: 2,
	})
}

// sendDietFromAI gets diet from AI and sends to user
func (h *LogicHandlers) sendDietFromAI(
	b *gotgbot.Bot, ctx *ext.Context,
	mealsAndSnacks model.MealsAndSnacksNumber,
) error {
	const errMsg = "Что-то пошло не так. Попробуйте еще раз"

	executionCtx := context.Background()

	// Check if user exists
	telegramID := ctx.EffectiveSender.Id()
	user, err := h.storage.GetUserByTelegramID(executionCtx, telegramID)
	if err != nil {
		// User not found
		if errors.Is(err, store.ErrNotFound) {
			ierr := sendMockedDiet(b, ctx, mealsAndSnacks)
			if ierr != nil {
				log.Printf("sendMockedDiet: %v", err)
				return endConversation(b, ctx, errMsg)
			}
			return h.SelectOrderProducts(b, ctx)
		}

		log.Printf("h.storage.GetUserByTelegramID: %v", err)
		return endConversation(b, ctx, errMsg)
	}

	// Check if user data filled
	if !user.IsFilledForGetDiet() {
		// TODO: send info about update data
		msg := "Не все параметры заполнены для составления диеты. Заполните и попробуйте снова"
		return endConversation(b, ctx, msg)
	}

	params := model.GetDietParams{
		Age:                  *user.Age,
		Weight:               *user.Weight,
		Height:               *user.Height,
		Gender:               *user.Gender,
		PhysicalActivity:     *user.PhysicalActivity,
		MealsAndSnacksNumber: mealsAndSnacks,
	}

	// Build prompt for AI
	prompt := buildPromptForGetDiet(params)

	// Get diet concurrently
	wg := &sync.WaitGroup{}
	var diet string
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		var ierr error
		diet, ierr = h.sendRequestToAI(ctx, prompt)
		if ierr != nil {
			log.Printf("h.sendRequestToAI: %v", ierr)
		}
	}(ctxutil.Detach(executionCtx))

	sendMessage(b, ctx, "Подождите немного. Составляю диету", nil)
	wg.Wait()

	// Create prompt log in db
	go func(ctx context.Context) {
		reqCtx, cancel := context.WithTimeout(ctx, createLogTimeout)
		defer cancel()

		dialogData := model.ChatBotDialogDataGetDietFromAI{Params: params}
		dialog := &model.ChatBotDialog{
			UserTelegramID: telegramID,
			Kind:           model.ChatBotDialogKindGetDietFromAI,
			Status:         model.ChatBotDialogStatusInProgress,
			DataJSON:       dialogData.ToJSON(),
		}
		ierr := h.storage.CreateChatBotDialog(reqCtx, dialog)
		if ierr != nil {
			log.Printf("h.storage.CreateAIAPILog: %v", ierr)
		}

		promptLog := &model.AIAPILog{
			Prompt:     prompt,
			Response:   ptrconv.Ptr(diet),
			UserID:     user.ID,
			SourceID:   dialog.ID,
			SourceType: model.AIAPILogsSourceTypeChatbotDialog,
		}
		ierr = h.storage.CreateAIAPILog(reqCtx, promptLog)
		if ierr != nil {
			log.Printf("h.storage.CreateAIAPILog: %v", ierr)
		}
	}(ctxutil.Detach(executionCtx))

	err = h.sendDiet(b, ctx, diet)
	if err != nil {
		log.Printf("h.sendDiet: %v", err)
		return endConversation(b, ctx, errMsg)
	}

	return h.SelectOrderProducts(b, ctx)
}

// SelectOrderProducts .
func (h *LogicHandlers) SelectOrderProducts(b *gotgbot.Bot, ctx *ext.Context) error {
	opts := &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "да", CallbackData: GetDietFromAISelectOrderProductsYes},
				{Text: "нет", CallbackData: GetDietFromAISelectOrderProductsNo},
			}},
		},
	}
	return replyInConversation(b, ctx, "Сделать заказ продуктов на Сбермаркете?", GetDietFromAISelectOrderProducts, opts)
}

// SelectOrderProductsYes .
func (h *LogicHandlers) SelectOrderProductsYes(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	// Send callback answer
	cb := ctx.Update.CallbackQuery
	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали: да",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали: да. Т.е. сделать заказ продуктов на Сбермаркете", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	// Create log in db
	go h.logNeedOrder(context.Background(), ctx.EffectiveSender.Id(), true)

	msg := fmt.Sprintf(`
На данный момент функционал находится в разработке.
Вы можете заказать продукты на Сбермаркете по этой ссылке:
%s`, sbermarketURL)
	return endConversation(b, ctx, msg)
}

// SelectOrderProductsNo .
func (h *LogicHandlers) SelectOrderProductsNo(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	// Send callback answer
	cb := ctx.Update.CallbackQuery
	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали: нет",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали: нет. Т.е. не делать заказ продуктов на Сбермаркете", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	// Create log in db
	go h.logNeedOrder(context.Background(), ctx.EffectiveSender.Id(), false)

	msg := "ОК, тогда в следующий раз"
	return endConversation(b, ctx, msg)
}

func (h *LogicHandlers) logNeedOrder(ctx context.Context, telegramID int64, needOrder bool) {
	logCtx, cancel := context.WithTimeout(ctx, createLogTimeout)
	defer cancel()

	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		logCtx, telegramID,
		model.ChatBotDialogKindGetDietFromAI,
		model.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)
		return
	}

	dialogData := &model.ChatBotDialogDataGetDietFromAI{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("GetDietFromAI.FromJSON: %v", err)
		return
	}

	dialogData.NeedOrder = ptrconv.Ptr(needOrder)
	dialog.DataJSON = dialogData.ToJSON()
	dialog.Status = model.ChatBotDialogStatusCompleted
	err = h.storage.UpdateChatBotDialog(logCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)
		return
	}
}

func (h *LogicHandlers) sendRequestToAI(ctx context.Context, prompt string) (string, error) {
	reqCtx, cancel := context.WithTimeout(ctx, sendRequestToAITimeout)
	defer cancel()

	diet, err := h.aiClient.SendRequest(reqCtx, prompt)
	if err != nil {
		return "", errors.Wrap(err, "h.aiClient.SendRequest")
	}

	return diet, nil
}

func (h *LogicHandlers) sendDiet(b *gotgbot.Bot, ctx *ext.Context, diet string) error {
	reader := strings.NewReader(diet)

	const caption = `
Вот диета для вас от ИИ (Искусственного Интеллекта).
Она носит рекомендательный характер.
Для получения более подходящей диеты обратитесь к диетологу:
/getdiet`
	err := sendDocument(b, ctx, reader, menuFileName, caption)
	if err != nil {
		return errors.Wrap(err, "sendDocument")
	}

	return nil
}

// buildPromptForGetDiet builds prompt for getting diet from AI
func buildPromptForGetDiet(params model.GetDietParams) string {
	// Get diet for interval fasting
	if params.MealsNumberPerDay == 2 {
		return strings.TrimSpace(fmt.Sprintf(
			promptTemplateGetIntervalDiet,
			params.Age,
			params.Height,
			params.Weight,
			params.Gender.DescriptionRu(),
			params.PhysicalActivity.DescriptionRu(),
		))
	}

	// Set snack times
	snacksNumber := "без перекусов"
	if params.SnacksNumberPerDay == 1 {
		snacksNumber = "с одним перекусом"
	} else if params.SnacksNumberPerDay == 2 {
		snacksNumber = "с двумя перекусами"
	}

	// Get usual diet
	return strings.TrimSpace(fmt.Sprintf(
		promptTemplateGetDiet,
		snacksNumber,
		params.Age,
		params.Height,
		params.Weight,
		params.Gender.DescriptionRu(),
		params.PhysicalActivity.DescriptionRu(),
	))
}
