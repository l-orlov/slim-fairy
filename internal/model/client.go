package model

import (
	"time"

	"github.com/google/uuid"
)

// ClientToRegister is model for registering client
type ClientToRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
	Age      int    `json:"age" binding:"required"`
	Weight   int    `json:"weight" binding:"required"`
	Height   int    `json:"height" binding:"required"`
	Gender   Gender `json:"gender" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ClientToSignIn is model for client sign-in
type ClientToSignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Client is model for client
type Client struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Age       int       `json:"age" db:"age"`
	Weight    int       `json:"weight" db:"weight"`
	Height    int       `json:"height" db:"height"`
	Gender    Gender    `json:"gender" db:"gender"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// DbTable returns DB table name
func (Client) DbTable() string {
	return "clients"
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
