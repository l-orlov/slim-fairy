package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

// ChatBotDialog is chatbot dialogs with user
type ChatBotDialog struct {
	ID             uuid.UUID           `db:"id"`
	UserTelegramID int64               `db:"user_telegram_id"`
	Kind           ChatBotDialogKind   `db:"kind"`
	Status         ChatBotDialogStatus `db:"status"`
	DataJSON       string              `db:"data"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// DbTable returns DB table name
func (ChatBotDialog) DbTable() string {
	return "chat_bot_dialogs"
}

// ChatBotDialogKind is type for ChatBotDialog kind
type ChatBotDialogKind string

// ChatBotDialogKind values
const (
	ChatBotDialogKindUserRegistration ChatBotDialogKind = "UserRegistration" // UserRegistration
	ChatBotDialogKindGetDietFromAI    ChatBotDialogKind = "GetDietFromAI"    // GetDietFromAI
)

// ChatBotDialogStatus is type for ChatBotDialog status
//
//go:generate stringer -type=ChatBotDialogStatus
type ChatBotDialogStatus int32

// ChatBotDialogStatus values
const (
	ChatBotDialogStatusInitial    = ChatBotDialogStatus(1) // Initial
	ChatBotDialogStatusInProgress = ChatBotDialogStatus(2) // InProgress
	ChatBotDialogStatusCanceled   = ChatBotDialogStatus(3) // Canceled
	ChatBotDialogStatusCompleted  = ChatBotDialogStatus(4) // Completed
)

type SelfMarshaller interface {
	Marshal() ([]byte, error)
}

type SelfUnmarshaller interface {
	Unmarshal(data []byte) error
}

// ChatBotDialogData is interface for storing dialog data
type ChatBotDialogData interface {
	driver.Valuer
	sql.Scanner
}

type ChatBotDialogDataUserRegistration struct {
	Name             *string                `json:"name,omitempty"`
	Age              *int                   `json:"age,omitempty"`
	Weight           *int                   `json:"weight,omitempty"`
	Height           *int                   `json:"height,omitempty"`
	Gender           *Gender                `json:"gender,omitempty"`
	PhysicalActivity *PhysicalActivityLevel `json:"physical_activity_level,omitempty"`
}

// ToJSON marshals data to JSON string
func (field *ChatBotDialogDataUserRegistration) ToJSON() string {
	data, err := json.Marshal(field)
	if err != nil {
		log.Printf("ChatBotDialogDataUserRegistration.ToJSON: %v", err)
	}

	return string(data)
}

// FromJSON unmarshals data form JSON string
func (field *ChatBotDialogDataUserRegistration) FromJSON(s string) error {
	return json.Unmarshal([]byte(s), field)
}

// IsFilled returns true if all fields are filled
func (field *ChatBotDialogDataUserRegistration) IsFilled() bool {
	return field.Name != nil && field.Age != nil && field.Weight != nil &&
		field.Height != nil && field.Gender != nil && field.PhysicalActivity != nil
}
