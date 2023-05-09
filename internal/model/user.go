package model

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// UserToRegister is model for registering user
type UserToRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
	Age      int    `json:"age" binding:"required"`
	Weight   int    `json:"weight" binding:"required"`
	Height   int    `json:"height" binding:"required"`
	Gender   Gender `json:"gender" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserToSignIn is model for user sign-in
type UserToSignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// User is model for user
type User struct {
	ID               uuid.UUID              `json:"id" db:"id"`
	Name             string                 `json:"name" db:"name"`
	Email            *string                `json:"email" db:"email"`
	Phone            *string                `json:"phone" db:"phone"`
	TelegramID       *int64                 `json:"telegram_id" db:"telegram_id"`
	Age              *int                   `json:"age" db:"age"`
	Weight           *int                   `json:"weight" db:"weight"`
	Height           *int                   `json:"height" db:"height"`
	Gender           *Gender                `json:"gender" db:"gender"`
	PhysicalActivity *PhysicalActivityLevel `json:"physical_activity_level" db:"physical_activity_level"`
	CreatedBy        UserCreatedBy          `json:"created_by" db:"created_by"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// DbTable returns DB table name
func (User) DbTable() string {
	return "users"
}

// IsFilledForGetDiet returns true if data is filled for diet
func (u User) IsFilledForGetDiet() bool {
	return u.Age != nil && u.Weight != nil &&
		u.Height != nil && u.Gender != nil && u.PhysicalActivity != nil
}

// Gender is type for person gender
//
//go:generate stringer -type=Gender
type Gender int32

// Gender values
const (
	GenderMale   = Gender(1) // Male
	GenderFemale = Gender(2) // Female
)

func (g Gender) DescriptionRu() string {
	switch g {
	case GenderMale:
		return "мужчина"
	case GenderFemale:
		return "женщина"
	default:
		log.Print("Gender.DescriptionRu: unknown value")
		return "невалидное значение"
	}
}

// PhysicalActivityLevel is type for person physical activity level
//
//go:generate stringer -type=PhysicalActivityLevel
type PhysicalActivityLevel int32

// PhysicalActivityLevel values
const (
	PhysicalActivityLevelLow    = PhysicalActivityLevel(1) // Low
	PhysicalActivityLevelMedium = PhysicalActivityLevel(2) // Medium
	PhysicalActivityLevelHigh   = PhysicalActivityLevel(3) // High
)

func (l PhysicalActivityLevel) DescriptionRu() string {
	switch l {
	case PhysicalActivityLevelLow:
		return "низкий"
	case PhysicalActivityLevelMedium:
		return "средний"
	case PhysicalActivityLevelHigh:
		return "высокий"
	default:
		log.Print("PhysicalActivityLevel.DescriptionRu: unknown value")
		return "невалидное значение"
	}
}

// UserCreatedBy is type for User CreatedBy field
type UserCreatedBy string

// UserCreatedBy values
const (
	UserCreatedByChatbot = UserCreatedBy("chatbot") // Chatbot
	UserCreatedByAPI     = UserCreatedBy("api")     // API
)
