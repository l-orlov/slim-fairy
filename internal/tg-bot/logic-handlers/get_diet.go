package logic_handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/l-orlov/slim-fairy/pkg/ctxutil"
	"github.com/l-orlov/slim-fairy/pkg/ptrconv"
	"github.com/pkg/errors"
)

// Mocked diets config
const (
	mocksNumber          = 13
	mockFilePathTemplate = "assets/mock_menus/%d.txt"
)

// Prompts templates
const (
	// Get usual diet
	promptTemplateGetDiet = `
Составь диету на неделю с тремя приемами пищи в день и %s.
Укажи список ингредиентов в конце.

Возраст: %d.
Рост: %d см.
Вес: %d кг.
Пол: %s.
Уровень физической активности: %s`
	// Get diet for interval fasting
	promptTemplateGetIntervalDiet = `
Составь диету для интервального голодания на неделю с двумя приемами пищи в день и без перекусов.
Укажи список ингредиентов в конце.

Возраст: %d.
Рост: %d см.
Вес: %d кг.
Пол: %s.
Уровень физической активности: %s`
)

const (
	menuFileName = "menu.txt"
)

// Timeouts
const (
	getDietFromAITimeout = 2 * time.Minute
	createLogTimeout     = 1 * time.Minute
)

// GetDietFromAI gets diet from AI and sends to user
func (h *LogicHandlers) GetDietFromAI(b *gotgbot.Bot, ctx *ext.Context) error {
	const errMsg = "Что-то пошло не так. Попробуйте еще раз"

	executionCtx := context.Background()

	// Check if user exists
	user, err := h.storage.GetUserByTelegramID(executionCtx, ctx.EffectiveSender.Id())
	if err != nil {
		// User not found
		if errors.Is(err, store.ErrNotFound) {
			ierr := sendMockedDiet(b, ctx)
			if ierr != nil {
				log.Printf("sendMockedDiet: %v", err)
				reply(b, ctx, errMsg)
				return nil
			}
			return nil
		}

		log.Printf("h.storage.GetUserByTelegramID: %v", err)
		reply(b, ctx, errMsg)
		return nil
	}

	// Check if user data filled
	if !user.IsFilledForGetDiet() {
		// TODO: send info about update data
		reply(b, ctx, "Не все параметры заполнены для составления диеты. Заполните и попробуйте снова")
		return nil
	}

	params := &model.GetDietParams{
		Age:              *user.Age,
		Weight:           *user.Weight,
		Height:           *user.Height,
		Gender:           *user.Gender,
		PhysicalActivity: *user.PhysicalActivity,
		// TODO: fix
		MealTimes:  3,
		SnackTimes: 1,
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
		diet, ierr = h.getDietFromAI(ctx, prompt)
		if ierr != nil {
			log.Printf("h.getDietFromAI: %v", ierr)
		}
	}(ctxutil.Detach(executionCtx))

	reply(b, ctx, "Подождите немного. Составляю диету")
	wg.Wait()

	// Create prompt log in db
	go func(ctx context.Context) {
		reqCtx, cancel := context.WithTimeout(ctx, createLogTimeout)
		defer cancel()

		promptLog := &model.AIAPILog{
			Prompt:   prompt,
			Response: ptrconv.Ptr(diet),
			UserID:   user.ID,
			// TODO: fill
			SourceID:   uuid.UUID{},
			SourceType: model.AIAPILogsSourceTypeChatbotDialog,
		}
		ierr := h.storage.CreateAIAPILog(reqCtx, promptLog)
		if ierr != nil {
			log.Printf("h.storage.CreateAIAPILog: %v", ierr)
		}
	}(ctxutil.Detach(executionCtx))

	err = h.sendDiet(b, ctx, diet)
	if err != nil {
		log.Printf("h.sendDiet: %v", err)
		reply(b, ctx, errMsg)
		return nil
	}

	return nil
}

func (h *LogicHandlers) getDietFromAI(ctx context.Context, prompt string) (string, error) {
	reqCtx, cancel := context.WithTimeout(ctx, getDietFromAITimeout)
	defer cancel()

	diet, err := h.aiClient.SendRequest(reqCtx, prompt)
	if err != nil {
		return "", errors.Wrap(err, "h.aiClient.SendRequest")
	}

	return diet, nil
}

func (h *LogicHandlers) sendDiet(b *gotgbot.Bot, ctx *ext.Context, diet string) error {
	reader := strings.NewReader(diet)

	const caption = "Вот диета для вас от ИИ (Искусственного Интеллекта)"
	err := sendDocument(b, ctx, reader, caption)
	if err != nil {
		return errors.Wrap(err, "sendDocument")
	}

	return nil
}

func sendMockedDiet(b *gotgbot.Bot, ctx *ext.Context) error {
	diet, err := getMockedDiet()
	if err != nil {
		return errors.Wrap(err, "getMockedDiet")
	}

	reader := strings.NewReader(diet)

	const caption = `
Вот пример диеты. Пройдите регистрацию, чтобы составил диету под ваши параметры:
/register
`
	err = sendDocument(b, ctx, reader, caption)
	if err != nil {
		return errors.Wrap(err, "sendDocument")
	}

	return nil
}

func sendDocument(
	b *gotgbot.Bot, ctx *ext.Context,
	file io.Reader, caption string) error {
	_, err := b.SendDocument(ctx.EffectiveChat.Id, gotgbot.NamedFile{
		File:     file,
		FileName: menuFileName,
	}, &gotgbot.SendDocumentOpts{
		Caption:          caption,
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	})
	if err != nil {
		return errors.Wrap(err, "b.SendDocument")
	}

	return nil
}

func getMockedDiet() (string, error) {
	fileNum := rand.Int31n(mocksNumber)
	path := fmt.Sprintf(mockFilePathTemplate, fileNum)

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "os.ReadFile")
	}

	return string(fileBytes), nil
}

// buildPromptForGetDiet builds prompt for getting diet from AI
func buildPromptForGetDiet(params *model.GetDietParams) string {
	// Get diet for interval fasting
	if params.MealTimes == 2 {
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
	snackTimes := "без перекусов"
	if params.SnackTimes == 1 {
		snackTimes = "с одним перекусом"
	} else if params.SnackTimes == 2 {
		snackTimes = "с двумя перекусами"
	}

	// Get usual diet
	return strings.TrimSpace(fmt.Sprintf(
		promptTemplateGetDiet,
		snackTimes,
		params.Age,
		params.Height,
		params.Weight,
		params.Gender.DescriptionRu(),
		params.PhysicalActivity.DescriptionRu(),
	))
}
