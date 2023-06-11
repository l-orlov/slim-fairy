package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/l-orlov/slim-fairy/bot/internal/ai-api-client"
	"github.com/l-orlov/slim-fairy/bot/internal/config"
	"github.com/l-orlov/slim-fairy/bot/internal/store"
	"github.com/l-orlov/slim-fairy/bot/internal/tg-bot"
)

func main() {
	ctx := context.Background()

	// Load config
	err := config.Load("configs/")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Get()

	// Connect to DB
	storage, err := store.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	// Create AI API client
	aiClient := ai_api_client.New(cfg.APIKey)

	// Create bot
	bot, err := tg_bot.New(aiClient, storage)
	if err != nil {
		log.Fatalf("tg_bot.New: %v", err)
	}

	// Run bot
	go bot.Run()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Printf("service shutting down")
	if err = bot.Stop(); err != nil {
		log.Printf("failed to stop bot: %v", err)
	}
}
