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

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/pkg/errors"
)

// Mocked diets config
const (
	mocksNumber          = 13
	mockFilePathTemplate = "assets/mock_menus/%d.txt"
)

const (
	menuFileName = "menu.txt"
)

// TODO: comment
func (h *LogicHandlers) GetDietFromAI(b *gotgbot.Bot, ctx *ext.Context) error {
	const errMsg = "Что-то пошло не так. Попробуйте еще раз"
	// если нет:
	// - вот пример диеты. пройдите регистрацию, чтобы составил диету под ваши параметры.
	// если есть:
	// - в горутине получение диеты
	// - Подождите немного. Составляю диету.
	// wg.Wait и потом ответ файликом

	// Check if user exists
	user, err := h.storage.GetUserByTelegramID(context.Background(), ctx.EffectiveSender.Id())
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
		MealTimes:  2,
		SnackTimes: 1,
	}

	// write request to db

	// Get diet concurrently
	wg := &sync.WaitGroup{}
	var diet string
	wg.Add(1)
	go func() {
		defer wg.Done()

		var dietErr error
		diet, dietErr = h.getDietFromAI(params)
		if dietErr != nil {
			log.Printf("h.getDietFromAI: %v", err)
		}
	}()

	reply(b, ctx, "Подождите немного. Составляю диету")
	wg.Wait()

	// write response to db

	err = h.sendDiet(b, ctx, diet)
	if err != nil {
		log.Printf("h.sendDiet: %v", err)
		reply(b, ctx, errMsg)
		return nil
	}

	return nil
}

func (h *LogicHandlers) getDietFromAI(params *model.GetDietParams) (diet string, err error) {
	executionCtx := context.Background()
	diet, err = h.dietGetter.GetDietByParams(executionCtx, params)
	if err != nil {
		return "", errors.Wrap(err, "h.dietGetter.GetDietByParams")
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
