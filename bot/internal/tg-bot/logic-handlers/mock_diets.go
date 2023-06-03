package logic_handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/l-orlov/slim-fairy/bot/internal/config"
	"github.com/l-orlov/slim-fairy/bot/internal/model"
	"github.com/pkg/errors"
)

// Mocked diets config
const (
	mockDietPathPrefix              = "assets/mock_diets/"
	mockDietMeals2PathPrefix        = "meals_2/"
	mockDietMeals3PathPrefix        = "meals_3/"
	mockDietMeals3Snacks0PathPrefix = "snacks_0/"
	mockDietMeals3Snacks1PathPrefix = "snacks_1/"
	mockDietMeals3Snacks2PathPrefix = "snacks_2/"
	mockDietFileNameTemplate        = "%d.txt"

	dietFileCaption = `
Вот пример диеты. Пройдите регистрацию, чтобы составил диету под ваши параметры:
/register`
)

func sendMockedDiet(
	b *gotgbot.Bot, ctx *ext.Context,
	mealsAndSnacks model.MealsAndSnacksNumber,
) error {
	diet, err := getMockedDiet(mealsAndSnacks)
	if err != nil {
		return errors.Wrap(err, "getMockedDiet")
	}

	reader := strings.NewReader(diet)

	err = sendDocument(b, ctx, reader, menuFileName, dietFileCaption)
	if err != nil {
		return errors.Wrap(err, "sendDocument")
	}

	return nil
}

func getMockedDiet(mealsAndSnacks model.MealsAndSnacksNumber) (string, error) {
	// Set path to mocks depending on params
	dirPath := mockDietPathPrefix
	if mealsAndSnacks.MealsNumberPerDay == 2 {
		dirPath += mockDietMeals2PathPrefix
	} else {
		dirPath += mockDietMeals3PathPrefix
		if mealsAndSnacks.SnacksNumberPerDay == 0 {
			dirPath += mockDietMeals3Snacks0PathPrefix
		} else if mealsAndSnacks.SnacksNumberPerDay == 1 {
			dirPath += mockDietMeals3Snacks1PathPrefix
		} else {
			dirPath += mockDietMeals3Snacks2PathPrefix
		}
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", errors.Wrap(err, "os.ReadDir")
	}

	// Choose random mock
	randomFileNum := config.RandomGenerator.Int31n(int32(len(files)))
	path := fmt.Sprintf(dirPath+mockDietFileNameTemplate, randomFileNum)

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "os.ReadFile")
	}

	return string(fileBytes), nil
}
