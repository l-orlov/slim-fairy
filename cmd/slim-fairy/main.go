package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/l-orlov/slim-fairy/internal/config"
	"github.com/l-orlov/slim-fairy/internal/handler"
	"github.com/l-orlov/slim-fairy/internal/service"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/l-orlov/slim-fairy/pkg/server"
	"github.com/pkg/errors"
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

	// Create service entity with business logic
	svc := service.New(storage)

	// Init HTTP handlers
	h := handler.New(storage, svc)

	// Start HTTP server
	srv := server.New(h, cfg.ServerAddress)
	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http server: %v", err)
		}
	}()

	log.Printf("service started %s", cfg.ServerAddress)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Printf("service shutting down")
	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shut down: %v", err)
	}
}
