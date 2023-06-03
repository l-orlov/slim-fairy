package model

import (
	"time"

	"github.com/google/uuid"
)

// Nutritionist is model for nutritionist
type Nutritionist struct {
	ID         uuid.UUID             `json:"id" db:"id"`
	Name       string                `json:"name" db:"name"`
	Email      *string               `json:"email" db:"email"`
	Phone      *string               `json:"phone" db:"phone"`
	TelegramID *int64                `json:"telegram_id" db:"telegram_id"`
	Age        *int                  `json:"age" db:"age"`
	Gender     *Gender               `json:"gender" db:"gender"`
	Info       *string               `json:"info" db:"info"`
	CreatedBy  NutritionistCreatedBy `json:"created_by" db:"created_by"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// DbTable returns DB table name
func (Nutritionist) DbTable() string {
	return "nutritionists"
}

// NutritionistCreatedBy is type for Nutritionist CreatedBy field
type NutritionistCreatedBy string

// NutritionistCreatedBy values
const (
	NutritionistCreatedByChatbot = NutritionistCreatedBy("chatbot") // Chatbot
	NutritionistCreatedByAPI     = NutritionistCreatedBy("api")     // API
)
