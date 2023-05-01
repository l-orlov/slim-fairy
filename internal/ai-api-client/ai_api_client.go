package ai_api_client

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/pkg/errors"
)

const (
	mocksNumber          = 7
	mockFilePathTemplate = "assets/mock_menus/%d.txt"
)

type (
	Client struct {
		// TODO: user to chatGPT
	}
)

func New() *Client {
	return &Client{}
}

func (c *Client) GetMenuByParams(ctx context.Context, params model.GetMenuParams) (string, error) {
	fileNum := rand.Int31n(mocksNumber)
	path := fmt.Sprintf(mockFilePathTemplate, fileNum)

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "os.ReadFile")
	}

	return string(fileBytes), nil
}
