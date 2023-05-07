package ai_api_client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/pkg/errors"
)

// ChatGPT client config
const (
	timeout     = 2 * time.Minute
	roleUser    = "user"
	maxTokens   = 3000
	temperature = 0.7
)

// todo: test
/*
Типа сделай меню для интервального голодания на 2 приема пищи без перекусов
*/

// ChatGPT prompts templates
const (
	// TODO: не работает, если 1 прием пищи и 1 перекус. Или если 2 прима пищи.
	// можно пофиксить через enum и захардкоженные подписи.
	promptTemplateGetDiet = `
Составь диету на неделю.
Приемов пищи: %d.
Перекусов: %d.
Без диетических ограничений.
Укажи список ингредиентов в конце.

Возраст: %d.
Рост: %d см.
Вес: %d кг.
Пол: %s.
Уровень физической активности: %s`
)

type (
	Client struct {
		client gpt3.Client
	}
)

func New(apiKey string) *Client {
	client := gpt3.NewClient(apiKey, gpt3.WithTimeout(timeout))

	return &Client{
		client: client,
	}
}

func (c *Client) GetDietByParams(ctx context.Context, params *model.GetDietParams) (string, error) {
	// Send request to ChatGPT API
	msg := buildPromptForGetDiet(params)
	output, err := c.sendRequest(ctx, msg)
	if err != nil {
		return "", errors.Wrap(err, "c.sendRequest")
	}

	return output, nil
}

func (c *Client) sendRequest(ctx context.Context, msg string) (resp string, err error) {
	outputBuilder := strings.Builder{}
	req := gpt3.ChatCompletionRequest{
		Messages: []gpt3.ChatCompletionRequestMessage{
			{
				Role:    roleUser,
				Content: msg,
			},
		},
		MaxTokens:   maxTokens,
		Temperature: gpt3.Float32Ptr(temperature),
	}
	err = c.client.ChatCompletionStream(ctx, req, func(resp *gpt3.ChatCompletionStreamResponse) {
		outputBuilder.WriteString(resp.Choices[0].Delta.Content)
	})
	if err != nil {
		return "", errors.Wrap(err, "c.client.ChatCompletionStream")
	}

	return outputBuilder.String(), nil
}

func buildPromptForGetDiet(params *model.GetDietParams) string {
	return fmt.Sprintf(
		promptTemplateGetDiet,
		params.MealTimes,
		params.SnackTimes,
		params.Age,
		params.Height,
		params.Weight,
		params.Gender.DescriptionRu(),
		params.PhysicalActivity.DescriptionRu(),
	)
}
