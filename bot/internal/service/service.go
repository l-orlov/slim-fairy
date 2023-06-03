package service

import (
	"github.com/l-orlov/slim-fairy/bot/internal/store"
)

// Service has methods with business logic
type Service struct {
	storage *store.Storage
}

// New creates new Service
func New(
	storage *store.Storage,
) *Service {
	return &Service{
		storage: storage,
	}
}
