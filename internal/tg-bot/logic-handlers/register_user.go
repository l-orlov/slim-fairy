package logic_handlers

import (
	"context"
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/l-orlov/slim-fairy/pkg/ptrconv"
	"github.com/pkg/errors"
)

// TODO:
// Escape string всюду. Как обезопаситься от sql injections?
// кнопки сделать, где это нужно.
// метод для обновления пользовательских данных.

// User registration handler states
const (
	RegisterName             = "name"
	RegisterAge              = "age"
	RegisterWeight           = "weight"
	RegisterHeight           = "height"
	RegisterGender           = "gender"
	RegisterPhysicalActivity = "physicalactivity"
	RegisterConfirm          = "сonfirm"
)

// userRegistrationStartInfo contains start info for user registration
const userRegistrationStartInfo = `
Чтобы точнее составить диету, нужны параметры (рост, вес и другие).
Я буду спрашивать по очереди каждый параметр.
Отвечайте мне в следующем сообщении.

Если решите прервать регистрацию, то используйте команду:
/cancelreg

Начнем. Как вас зовут?`

// StartUserRegistration starts user registration
func (h *LogicHandlers) StartUserRegistration(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	var (
		replyMsg            string
		needEndConversation bool
	)
	defer func() {
		// End if needed
		if needEndConversation {
			nextState = endConversation(b, ctx, replyMsg)
			return
		}

		// Reply to user
		nextState = replyInConversation(b, ctx, replyMsg, RegisterName)
	}()

	// Check if user exists
	user, err := h.storage.GetUserByTelegramID(context.Background(), ctx.EffectiveSender.Id())
	if err == nil {
		// TODO: подсказка о том, как обновить данные, если это нужно
		replyMsg = fmt.Sprintf("%s, вы уже прошли регистрацию", user.Name)
		needEndConversation = true
		return nil
	}
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		log.Printf("h.storage.GetUserByTelegramID: %v", err)

		// Error -> try again
		replyMsg = "Что-то пошло не так. Попробуйте еще раз"
		return nil
	}

	// Success -> go to next state
	replyMsg = userRegistrationStartInfo

	return nil
}

// CancelUserRegistration cancels user registration
func (h *LogicHandlers) CancelUserRegistration(b *gotgbot.Bot, ctx *ext.Context) error {
	return endConversation(b, ctx, "ОК, тогда в следующий раз")
}

// RegisterUserName registers user name
func (h *LogicHandlers) RegisterUserName(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text

	var (
		replyMsg string
		success  bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, replyMsg, success, RegisterName, RegisterAge)
	}()

	// Validate name
	const validLen = 100
	if len(input) > validLen {
		// Not valid -> try again
		replyMsg = fmt.Sprintf("Слишком длинное имя. Сократите до %d символов", validLen)
		return nil
	}

	name := html.EscapeString(input)

	// Create dialog data
	dialogData := &model.ChatBotDialogDataUserRegistration{
		Name: ptrconv.Ptr(name),
	}
	dialog := &model.ChatBotDialog{
		UserTelegramID: ctx.EffectiveSender.Id(),
		Kind:           model.ChatBotDialogKindUserRegistration,
		Status:         model.ChatBotDialogStatusInProgress,
		DataJSON:       dialogData.ToJSON(),
	}
	err := h.storage.CreateChatBotDialog(context.Background(), dialog)
	if err != nil {
		log.Printf("h.storage.CreateChatBotDialog: %v", err)

		// Error -> try again
		replyMsg = "Что-то пошло не так. Введите имя еще раз"
		return nil
	}

	// Success -> go to next state
	replyMsg = fmt.Sprintf("Рад знакомству, %s!\nСколько вам лет?", name)
	success = true

	return nil
}

// RegisterUserAge registers user age
func (h *LogicHandlers) RegisterUserAge(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text

	var (
		replyMsg string
		success  bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, replyMsg, success, RegisterAge, RegisterWeight)
	}()

	ageNumber, err := strconv.Atoi(input)
	if err != nil {
		// Not valid -> try again
		replyMsg = "Нужно число. Попробуйте еще раз"
		return nil
	}
	if ageNumber < 0 || ageNumber > 150 {
		// Not valid -> try again
		replyMsg = "Число не подходит для возраста. Попробуйте еще раз"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите возраст еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model.ChatBotDialogKindUserRegistration,
		model.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData := &model.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData.Age = ptrconv.Ptr(ageNumber)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	// Success -> go to next state
	replyMsg = "Прекрасный возраст.\nКакой у вас вес кг?"
	success = true

	return nil
}

// RegisterUserWeight registers user weight
func (h *LogicHandlers) RegisterUserWeight(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text

	var (
		replyMsg string
		success  bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, replyMsg, success, RegisterWeight, RegisterHeight)
	}()

	weight, err := strconv.Atoi(input)
	if err != nil {
		// Not valid -> try again
		replyMsg = "Нужно число. Попробуйте еще раз"
		return nil
	}
	if weight < 0 || weight > 300 {
		// Not valid -> try again
		replyMsg = "Число не подходит для веса. Попробуйте еще раз"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите вес еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model.ChatBotDialogKindUserRegistration,
		model.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData := &model.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData.Weight = ptrconv.Ptr(weight)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	// Success -> go to next state
	replyMsg = "Хорошо.\nКакой у вас рост см?"
	success = true

	return nil
}

// RegisterUserHeight registers user height
func (h *LogicHandlers) RegisterUserHeight(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text

	var (
		replyMsg string
		success  bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, replyMsg, success, RegisterHeight, RegisterGender)
	}()

	height, err := strconv.Atoi(input)
	if err != nil {
		// Not valid -> try again
		replyMsg = "Нужно число. Попробуйте еще раз"
		return nil
	}
	if height < 0 || height > 250 {
		// Not valid -> try again
		replyMsg = "Число не подходит для роста. Попробуйте еще раз"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите рост еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model.ChatBotDialogKindUserRegistration,
		model.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData := &model.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData.Height = ptrconv.Ptr(height)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	// Success -> go to next state
	replyMsg = "ОК.\nКакой у вас пол?\nНапишите: м (мужчина) или ж (женщина)"
	success = true

	return nil
}

// RegisterUserGender registers user Gender
func (h *LogicHandlers) RegisterUserGender(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text
	input = strings.ToLower(input)

	var (
		replyMsg string
		success  bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, replyMsg, success, RegisterGender, RegisterPhysicalActivity)
	}()

	if input != "м" && input != "ж" {
		// Not valid -> try again
		replyMsg = "Попробуйте еще раз. Напишите: м или ж"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите пол еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model.ChatBotDialogKindUserRegistration,
		model.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData := &model.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData.Gender = ptrconv.Ptr(model.GenderMan)
	if input == "ж" {
		dialogData.Gender = ptrconv.Ptr(model.GenderWoman)
	}
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	// Success -> go to next state
	replyMsg = `
Принято.
Теперь какой у вас уровень физической активности?
Напишите: н (низкий) или с (средний) или в (высокий)`
	success = true

	return nil
}

// RegisterUserPhysicalActivity registers user physical activity
func (h *LogicHandlers) RegisterUserPhysicalActivity(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text
	input = strings.ToLower(input)

	var (
		replyMsg            string
		success             bool
		needEndConversation bool
	)
	defer func() {
		// End if needed
		if needEndConversation {
			nextState = endConversation(b, ctx, replyMsg)
			return
		}

		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, replyMsg, success, RegisterPhysicalActivity, RegisterConfirm)
	}()

	if input != "н" && input != "с" && input != "в" {
		// Not valid -> try again
		replyMsg = "Попробуйте еще раз. Напишите: н или с или в"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите пол еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model.ChatBotDialogKindUserRegistration,
		model.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData := &model.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData.PhysicalActivity = ptrconv.Ptr(model.PhysicalActivityLevelLow)
	if input == "с" {
		dialogData.PhysicalActivity = ptrconv.Ptr(model.PhysicalActivityLevelMedium)
	} else if input == "в" {
		dialogData.PhysicalActivity = ptrconv.Ptr(model.PhysicalActivityLevelHigh)
	}
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	// Check all data is filled
	if !dialogData.IsFilled() {
		log.Print("h.RegisterUserPhysicalActivity: not all data is filled")

		// Error -> try again
		replyMsg = `
Не все данные заполнены. Попробуйте пройти регистрацию снова:
/register`
		needEndConversation = true
		return nil
	}

	// Success -> go to next state
	replyMsg = fmt.Sprintf(`
Принято.

Имя: %s
Возраст: %d
Вес: %d
Рост: %d
Пол: %s
Уровень физической активности: %s

Подтвердите, что данные верны.
Напишите: да (верны) или нет (не верны)`,
		*dialogData.Name, *dialogData.Age, *dialogData.Weight, *dialogData.Height,
		dialogData.Gender.DescriptionRu(), dialogData.PhysicalActivity.DescriptionRu())
	success = true

	return nil
}

// ConfirmUserRegistration confirms user registration
func (h *LogicHandlers) ConfirmUserRegistration(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text
	input = strings.ToLower(input)

	var (
		replyMsg            string
		needEndConversation bool
	)
	defer func() {
		// End if needed
		if needEndConversation {
			nextState = endConversation(b, ctx, replyMsg)
			return
		}

		// Reply to user
		nextState = replyInConversation(b, ctx, replyMsg, RegisterConfirm)
	}()

	if input != "да" && input != "нет" {
		// Not valid -> try again
		replyMsg = "Попробуйте еще раз. Напишите: да или нет"
		return nil
	}

	if input == "нет" {
		replyMsg = `
ОК. Похоже какие-то данные неверны.
Попробуйте пройти регистрацию снова:
/register`
		needEndConversation = true
		return nil
	}

	// Use dialog data for creating user
	const errMsg = "Что-то пошло не так. Напишите еще раз: да или нет"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model.ChatBotDialogKindUserRegistration,
		model.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	dialogData := &model.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	user := buildUserFromDialogData(dialog.UserTelegramID, dialogData)
	txErr := h.storage.WithTransaction(dbCtx, func(tx store.Tx) error {
		err = h.storage.CreateUserTx(dbCtx, tx, user)
		if err != nil {
			return errors.Wrap(err, "h.storage.CreateUserTx")
		}

		err = h.storage.UpdateChatBotDialogStatusTx(dbCtx, tx, model.ChatBotDialogStatusCompleted, dialog.ID)
		if err != nil {
			return errors.Wrap(err, "h.storage.UpdateChatBotDialogStatusTx")
		}

		return nil
	})
	if txErr != nil {
		log.Printf("h.storage.WithTransaction: %v", err)

		// Error -> try again
		replyMsg = errMsg
		return nil
	}

	replyMsg = `
Вы успешно прошли регистрацию!
Уже учел ваши параметры для составления диеты.`
	needEndConversation = true
	return nil
}

func replyWithStatesInConversation(
	b *gotgbot.Bot, ctx *ext.Context,
	replyMsg string, success bool,
	failureStateName, successStateName string,
) (nextState error) {
	// If not success -> try again
	if !success {
		return replyInConversation(b, ctx, replyMsg, failureStateName)
	}

	// If success -> go to next state
	return replyInConversation(b, ctx, replyMsg, successStateName)
}

func replyInConversation(
	b *gotgbot.Bot, ctx *ext.Context,
	replyMsg string, nextStateName string,
) (nextState error) {
	reply(b, ctx, replyMsg)
	return handlers.NextConversationState(nextStateName)
}

func endConversation(b *gotgbot.Bot, ctx *ext.Context, replyMsg string) (nextState error) {
	reply(b, ctx, replyMsg)
	return handlers.EndConversation()
}

func reply(b *gotgbot.Bot, ctx *ext.Context, replyMsg string) {
	_, err := ctx.EffectiveMessage.Reply(b, replyMsg, &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		log.Printf("ctx.EffectiveMessage.Reply: %v", err)
	}
}

func buildUserFromDialogData(telegramID int64, data *model.ChatBotDialogDataUserRegistration) *model.User {
	user := &model.User{
		TelegramID:       ptrconv.Ptr(telegramID),
		Age:              data.Age,
		Weight:           data.Weight,
		Height:           data.Height,
		Gender:           data.Gender,
		PhysicalActivity: data.PhysicalActivity,
		CreatedBy:        model.UserCreatedByChatbot,
	}
	if data.Name != nil {
		user.Name = *data.Name
	}

	return user
}
