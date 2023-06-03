package logic_handlers

import (
	"context"
	"fmt"
	"html"
	"log"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	model2 "github.com/l-orlov/slim-fairy/bot/internal/model"
	store2 "github.com/l-orlov/slim-fairy/bot/internal/store"
	"github.com/l-orlov/slim-fairy/bot/pkg/ptrconv"
	"github.com/pkg/errors"
)

// TODO:
// Escape string всюду. Как обезопаситься от sql injections?
// метод для обновления пользовательских данных.

// Handler states for user registration
const (
	RegisterName             = "name"
	RegisterAge              = "age"
	RegisterWeight           = "weight"
	RegisterHeight           = "height"
	RegisterGender           = "gender"
	RegisterPhysicalActivity = "physicalactivity"
	RegisterConfirm          = "сonfirm"
)

// Callback keys for getting diet from AI
const (
	RegisterGenderCbMale             = "gender_male"
	RegisterGenderCbFemale           = "gender_female"
	RegisterPhysicalActivityCbLow    = "physicalactivity_low"
	RegisterPhysicalActivityCbMedium = "physicalactivity_medium"
	RegisterPhysicalActivityCbHigh   = "physicalactivity_high"
	RegisterConfirmCbYes             = "сonfirm_yes"
	RegisterConfirmCbNo              = "сonfirm_no"
)

// userRegistrationStartInfo contains start info for user registration
const userRegistrationStartInfo = `
Чтобы точнее составить диету, нужны параметры (рост, вес и другие).
Я буду спрашивать по очереди каждый параметр.
Отвечайте мне в следующем сообщении.

Если решите прервать регистрацию, то используйте команду:
/cancel

Начнем. Как вас зовут?`

// StartUserRegistration starts user registration
func (h *LogicHandlers) StartUserRegistration(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	var (
		msg                 string
		needEndConversation bool
	)
	defer func() {
		// End if needed
		if needEndConversation {
			nextState = endConversation(b, ctx, msg)
			return
		}

		// Reply to user
		nextState = replyInConversation(b, ctx, msg, RegisterName, nil)
	}()

	// Check if user exists
	user, err := h.storage.GetUserByTelegramID(context.Background(), ctx.EffectiveSender.Id())
	if err == nil {
		// TODO: подсказка о том, как обновить данные, если это нужно
		msg = fmt.Sprintf("%s, вы уже прошли регистрацию", user.Name)
		needEndConversation = true
		return nil
	}
	if err != nil && !errors.Is(err, store2.ErrNotFound) {
		log.Printf("h.storage.GetUserByTelegramID: %v", err)

		// Error -> try again
		msg = "Что-то пошло не так. Попробуйте еще раз"
		return nil
	}

	// Success -> go to next state
	msg = userRegistrationStartInfo

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
		msg     string
		success bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, msg, success, RegisterName, RegisterAge, nil)
	}()

	// Validate name
	const validLen = 100
	if len(input) > validLen {
		// Not valid -> try again
		msg = fmt.Sprintf("Слишком длинное имя. Сократите до %d символов", validLen)
		return nil
	}

	name := html.EscapeString(input)

	// Create dialog data
	dialogData := &model2.ChatBotDialogDataUserRegistration{
		Name: ptrconv.Ptr(name),
	}
	dialog := &model2.ChatBotDialog{
		UserTelegramID: ctx.EffectiveSender.Id(),
		Kind:           model2.ChatBotDialogKindUserRegistration,
		Status:         model2.ChatBotDialogStatusInProgress,
		DataJSON:       dialogData.ToJSON(),
	}
	err := h.storage.CreateChatBotDialog(context.Background(), dialog)
	if err != nil {
		log.Printf("h.storage.CreateChatBotDialog: %v", err)

		// Error -> try again
		msg = "Что-то пошло не так. Введите имя еще раз"
		return nil
	}

	// Success -> go to next state
	msg = fmt.Sprintf("Рад знакомству, %s!\nСколько вам лет?", name)
	success = true

	return nil
}

// RegisterUserAge registers user age
func (h *LogicHandlers) RegisterUserAge(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text

	var (
		msg     string
		success bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, msg, success, RegisterAge, RegisterWeight, nil)
	}()

	ageNumber, err := strconv.Atoi(input)
	if err != nil {
		// Not valid -> try again
		msg = "Нужно число. Попробуйте еще раз"
		return nil
	}
	if ageNumber < 0 || ageNumber > 150 {
		// Not valid -> try again
		msg = "Число не подходит для возраста. Попробуйте еще раз"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите возраст еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model2.ChatBotDialogKindUserRegistration,
		model2.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData := &model2.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData.Age = ptrconv.Ptr(ageNumber)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	// Success -> go to next state
	msg = "Прекрасный возраст.\nКакой у вас вес кг?"
	success = true

	return nil
}

// RegisterUserWeight registers user weight
func (h *LogicHandlers) RegisterUserWeight(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text

	var (
		msg     string
		success bool
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, msg, success, RegisterWeight, RegisterHeight, nil)
	}()

	weight, err := strconv.Atoi(input)
	if err != nil {
		// Not valid -> try again
		msg = "Нужно число. Попробуйте еще раз"
		return nil
	}
	if weight < 0 || weight > 300 {
		// Not valid -> try again
		msg = "Число не подходит для веса. Попробуйте еще раз"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите вес еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model2.ChatBotDialogKindUserRegistration,
		model2.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData := &model2.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData.Weight = ptrconv.Ptr(weight)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	// Success -> go to next state
	msg = "Хорошо.\nКакой у вас рост см?"
	success = true

	return nil
}

// RegisterUserHeight registers user height
func (h *LogicHandlers) RegisterUserHeight(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	input := ctx.EffectiveMessage.Text

	var (
		msg     string
		success bool
		opts    *gotgbot.SendMessageOpts
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, msg, success, RegisterHeight, RegisterGender, opts)
	}()

	height, err := strconv.Atoi(input)
	if err != nil {
		// Not valid -> try again
		msg = "Нужно число. Попробуйте еще раз"
		return nil
	}
	if height < 0 || height > 250 {
		// Not valid -> try again
		msg = "Число не подходит для роста. Попробуйте еще раз"
		return nil
	}

	// Update dialog data
	const errMsg = "Что-то пошло не так. Введите рост еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model2.ChatBotDialogKindUserRegistration,
		model2.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData := &model2.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData.Height = ptrconv.Ptr(height)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	// Success -> go to next state
	msg = "ОК.\nВаш пол? Выберите ниже"
	success = true
	opts = &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "мужчина", CallbackData: RegisterGenderCbMale},
				{Text: "женщина", CallbackData: RegisterGenderCbFemale},
			}},
		},
	}

	return nil
}

// RegisterUserGenderMale registers user gender with Male value
func (h *LogicHandlers) RegisterUserGenderMale(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	return h.registerUserGender(b, ctx, model2.GenderMale)
}

// RegisterUserGenderFemale registers user gender with Female value
func (h *LogicHandlers) RegisterUserGenderFemale(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	return h.registerUserGender(b, ctx, model2.GenderFemale)
}

// registerUserGender registers user gender
func (h *LogicHandlers) registerUserGender(
	b *gotgbot.Bot, ctx *ext.Context, gender model2.Gender,
) (nextState error) {
	var (
		msg     string
		success bool
		opts    *gotgbot.SendMessageOpts
	)
	defer func() {
		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, msg, success, RegisterGender, RegisterPhysicalActivity, opts)
	}()

	// Update dialog data
	const errMsg = "Что-то пошло не так. Выберите ваш пол еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model2.ChatBotDialogKindUserRegistration,
		model2.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData := &model2.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData.Gender = ptrconv.Ptr(gender)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	// Send callback answer
	genderDesription := gender.DescriptionRu()
	cb := ctx.Update.CallbackQuery
	_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали: " + genderDesription,
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали пол: "+genderDesription, nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	// Success -> go to next state
	msg = "Принято.\nВаш уровень физической активности? Выберите ниже"
	success = true
	opts = &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "низкий", CallbackData: RegisterPhysicalActivityCbLow},
				{Text: "средний", CallbackData: RegisterPhysicalActivityCbMedium},
				{Text: "высокий", CallbackData: RegisterPhysicalActivityCbHigh},
			}},
		},
	}

	return nil
}

// RegisterUserPhysicalActivityLow registers user physical activity with Low value
func (h *LogicHandlers) RegisterUserPhysicalActivityLow(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	return h.registerUserPhysicalActivity(b, ctx, model2.PhysicalActivityLevelLow)
}

// RegisterUserPhysicalActivityMedium registers user physical activity with Medium value
func (h *LogicHandlers) RegisterUserPhysicalActivityMedium(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	return h.registerUserPhysicalActivity(b, ctx, model2.PhysicalActivityLevelMedium)
}

// RegisterUserPhysicalActivityHigh registers user physical activity with High value
func (h *LogicHandlers) RegisterUserPhysicalActivityHigh(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	return h.registerUserPhysicalActivity(b, ctx, model2.PhysicalActivityLevelHigh)
}

// registerUserPhysicalActivity registers user physical activity
func (h *LogicHandlers) registerUserPhysicalActivity(
	b *gotgbot.Bot, ctx *ext.Context, activityLevel model2.PhysicalActivityLevel,
) (nextState error) {
	var (
		msg                 string
		success             bool
		needEndConversation bool
		opts                *gotgbot.SendMessageOpts
	)
	defer func() {
		// End if needed
		if needEndConversation {
			nextState = endConversation(b, ctx, msg)
			return
		}

		// Reply to user
		nextState = replyWithStatesInConversation(b, ctx, msg, success, RegisterPhysicalActivity, RegisterConfirm, opts)
	}()

	// Update dialog data
	const errMsg = "Что-то пошло не так. Выберите пол еще раз"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model2.ChatBotDialogKindUserRegistration,
		model2.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData := &model2.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData.PhysicalActivity = ptrconv.Ptr(activityLevel)
	dialog.DataJSON = dialogData.ToJSON()
	err = h.storage.UpdateChatBotDialog(dbCtx, dialog)
	if err != nil {
		log.Printf("h.storage.UpdateChatBotDialog: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	// Send callback answer
	activityDesription := activityLevel.DescriptionRu()
	cb := ctx.Update.CallbackQuery
	_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали: " + activityDesription,
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали уровень физической активности: "+activityDesription, nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	// Check all data is filled
	if !dialogData.IsFilled() {
		log.Print("h.RegisterUserPhysicalActivity: not all data is filled")

		// Error -> try again
		msg = `
Не все данные заполнены. Попробуйте пройти регистрацию снова:
/register`
		needEndConversation = true
		return nil
	}

	// Success -> go to next state
	msg = fmt.Sprintf(`
Принято.

Имя: %s
Возраст: %d
Вес: %d
Рост: %d
Пол: %s
Уровень физической активности: %s

Подтвердите, что данные верны.
Выберите: да (верны) или нет (не верны)`,
		*dialogData.Name, *dialogData.Age, *dialogData.Weight, *dialogData.Height,
		dialogData.Gender.DescriptionRu(), dialogData.PhysicalActivity.DescriptionRu())
	success = true
	opts = &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "да", CallbackData: RegisterConfirmCbYes},
				{Text: "нет", CallbackData: RegisterConfirmCbNo},
			}},
		},
	}

	return nil
}

// ConfirmUserRegistrationYes confirms user registration with Yes value
func (h *LogicHandlers) ConfirmUserRegistrationYes(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	var (
		msg                 string
		needEndConversation bool
	)
	defer func() {
		// End if needed
		if needEndConversation {
			nextState = endConversation(b, ctx, msg)
			return
		}

		// Reply to user
		nextState = replyInConversation(b, ctx, msg, RegisterConfirm, nil)
	}()

	// Use dialog data for creating user
	const errMsg = "Что-то пошло не так. Выберите еще раз: да или нет"
	dbCtx := context.Background()
	dialog, err := h.storage.GetChatBotDialogByKeyFields(
		dbCtx,
		ctx.EffectiveSender.Id(),
		model2.ChatBotDialogKindUserRegistration,
		model2.ChatBotDialogStatusInProgress,
	)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	dialogData := &model2.ChatBotDialogDataUserRegistration{}
	err = dialogData.FromJSON(dialog.DataJSON)
	if err != nil {
		log.Printf("h.storage.GetChatBotDialogByKeyFields: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	user := buildUserFromDialogData(dialog.UserTelegramID, dialogData)
	txErr := h.storage.WithTransaction(dbCtx, func(tx store2.Tx) error {
		err = h.storage.CreateUserTx(dbCtx, tx, user)
		if err != nil {
			return errors.Wrap(err, "h.storage.CreateUserTx")
		}

		err = h.storage.UpdateChatBotDialogStatusTx(dbCtx, tx, model2.ChatBotDialogStatusCompleted, dialog.ID)
		if err != nil {
			return errors.Wrap(err, "h.storage.UpdateChatBotDialogStatusTx")
		}

		return nil
	})
	if txErr != nil {
		log.Printf("h.storage.WithTransaction: %v", err)

		// Error -> try again
		msg = errMsg
		return nil
	}

	// Send callback answer
	cb := ctx.Update.CallbackQuery
	_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали: да",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали: да. Введенные данные верны", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	msg = `
Вы успешно прошли регистрацию!
Уже учел ваши параметры для составления диеты.
Получить диету от ИИ:
/getdietfromai`
	needEndConversation = true
	return nil
}

// ConfirmUserRegistrationNo confirms user registration with No value
func (h *LogicHandlers) ConfirmUserRegistrationNo(b *gotgbot.Bot, ctx *ext.Context) (nextState error) {
	// Send callback answer
	cb := ctx.Update.CallbackQuery
	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Вы выбрали: нет",
	})
	if err != nil {
		log.Printf("cb.Answer: %v", err)
	}

	_, _, err = cb.Message.EditText(b, "Вы выбрали: нет. Введенные данные неверны", nil)
	if err != nil {
		log.Printf("cb.Message.EditText: %v", err)
	}

	msg := `
ОК. Данные неверны.
Попробуйте пройти регистрацию снова:
/register`
	return endConversation(b, ctx, msg)
}

func buildUserFromDialogData(telegramID int64, data *model2.ChatBotDialogDataUserRegistration) *model2.User {
	user := &model2.User{
		TelegramID:       ptrconv.Ptr(telegramID),
		Age:              data.Age,
		Weight:           data.Weight,
		Height:           data.Height,
		Gender:           data.Gender,
		PhysicalActivity: data.PhysicalActivity,
		CreatedBy:        model2.UserCreatedByChatbot,
	}
	if data.Name != nil {
		user.Name = *data.Name
	}

	return user
}
