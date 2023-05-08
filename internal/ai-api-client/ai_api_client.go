package ai_api_client

import (
	"context"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/pkg/errors"
)

// ChatGPT client config
const (
	timeout     = 2 * time.Minute
	roleUser    = "user"
	maxTokens   = 3000
	temperature = 0.7
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

// SendRequest sends request to ChatGPT API and returns response
func (c *Client) SendRequest(ctx context.Context, prompt string) (string, error) {
	output, err := c.sendRequest(ctx, prompt)
	if err != nil {
		return "", errors.Wrap(err, "c.sendRequest")
	}

	return output, nil
}

func (c *Client) sendRequest(ctx context.Context, prompt string) (resp string, err error) {
	outputBuilder := strings.Builder{}
	req := gpt3.ChatCompletionRequest{
		Messages: []gpt3.ChatCompletionRequestMessage{
			{
				Role:    roleUser,
				Content: prompt,
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
