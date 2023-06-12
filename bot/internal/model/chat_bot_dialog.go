package model

import (
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

// ChatBotDialogDataUserRegistration .
type ChatBotDialogDataUserRegistration struct {
	Name             *string                `json:"name,omitempty"`
	Age              *int                   `json:"age,omitempty"`
	Weight           *int                   `json:"weight,omitempty"`
	Height           *int                   `json:"height,omitempty"`
	Gender           *Gender                `json:"gender,omitempty"`
	PhysicalActivity *PhysicalActivityLevel `json:"physical_activity_level,omitempty"`
}

// ToJSON marshals data to JSON string
func (data *ChatBotDialogDataUserRegistration) ToJSON() string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ChatBotDialogDataUserRegistration.ToJSON: %v", err)
	}

	return string(jsonBytes)
}

// FromJSON unmarshals data form JSON string
func (data *ChatBotDialogDataUserRegistration) FromJSON(s string) error {
	return json.Unmarshal([]byte(s), data)
}

// IsFilled returns true if all fields are filled
func (data *ChatBotDialogDataUserRegistration) IsFilled() bool {
	return data.Name != nil && data.Age != nil && data.Weight != nil &&
		data.Height != nil && data.Gender != nil && data.PhysicalActivity != nil
}

// ChatBotDialogDataGetDietFromAI .
type ChatBotDialogDataGetDietFromAI struct {
	Params    GetDietParams `json:"params"`
	NeedOrder *bool         `json:"need_order,omitempty"`
}

// ToJSON marshals data to JSON string
func (data *ChatBotDialogDataGetDietFromAI) ToJSON() string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ChatBotDialogDataUserRegistration.ToJSON: %v", err)
	}

	return string(jsonBytes)
}

// FromJSON unmarshals data form JSON string
func (data *ChatBotDialogDataGetDietFromAI) FromJSON(s string) error {
	return json.Unmarshal([]byte(s), data)
}
