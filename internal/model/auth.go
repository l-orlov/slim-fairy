package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthData is model for authentication data
type AuthData struct {
	SourceID   uuid.UUID          `db:"source_id"`
	SourceType AuthDataSourceType `db:"source_type"`
	Password   string             `db:"password"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// DbTable returns DB table name
func (AuthData) DbTable() string {
	return "auth_data"
}

// AuthDataSourceType is source type for authentication
type AuthDataSourceType string

// AuthDataSourceType values
const (
	AuthDataSourceTypeUser         = AuthDataSourceType("user")         // User
	AuthDataSourceTypeNutritionist = AuthDataSourceType("nutritionist") // Nutritionist
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
