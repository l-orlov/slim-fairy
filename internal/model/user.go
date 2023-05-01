package model

import (
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
	TelegramID       *string                `json:"telegram_id"`
	Age              *int                   `json:"age" db:"age"`
	Weight           *int                   `json:"weight" db:"weight"`
	Height           *int                   `json:"height" db:"height"`
	Gender           *Gender                `json:"gender" db:"gender"`
	PhysicalActivity *PhysicalActivityLevel `json:"physical_activity_level" db:"physical_activity_level"`
	CreatedAt        time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time              `json:"updatedAt" db:"updated_at"`
}

// DbTable returns DB table name
func (User) DbTable() string {
	return "users"
}

// Gender is type for person gender
//
//go:generate stringer -type=Gender
type Gender int32

// Gender values
const (
	GenderMan   = Gender(1) // Man
	GenderWoman = Gender(2) // Woman
)

// PhysicalActivityLevel is type for person physical activity level
//
//go:generate stringer -type=PhysicalActivityLevel
type PhysicalActivityLevel int32

// PhysicalActivityLevel values
const (
	PhysicalActivityLevelLow    = Gender(1) // Low
	PhysicalActivityLevelMedium = Gender(2) // Medium
	PhysicalActivityLevelHigh   = Gender(3) // High
)
