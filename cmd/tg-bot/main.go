package main

import (
	"context"
	"log"

	"github.com/l-orlov/slim-fairy/internal/ai-api-client"
	"github.com/l-orlov/slim-fairy/internal/config"
	"github.com/l-orlov/slim-fairy/internal/tg-bot"
)

func main() {
	ctx := context.Background()

	// Load config
	err := config.Load("configs/")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Get()

	// Create AI API user
	aiAPIUser := ai_api_client.New()

	// Create bot
	bot, err := tg_bot.New(cfg.Token, aiAPIUser)
	if err != nil {
		log.Fatalf("tg_bot.New: %v", err)
	}

	// Run bot
	bot.Run()

	_ = ctx
}
