package model

import (
	"time"

	"github.com/google/uuid"
)

// AIAPILog is model for AI API Logs
type AIAPILog struct {
	ID         uuid.UUID           `db:"id"`
	Prompt     string              `db:"prompt"`
	Response   *string             `db:"response"`
	UserID     uuid.UUID           `db:"user_id"`
	SourceID   uuid.UUID           `db:"source_id"`
	SourceType AIAPILogsSourceType `db:"source_type"`
	CreatedAt  time.Time           `db:"created_at"`
	UpdatedAt  time.Time           `db:"updated_at"`
}

// DbTable returns DB table name
func (AIAPILog) DbTable() string {
	return "ai_api_logs"
}

// AIAPILogsSourceType is type for AIAPILogs SourceType field
type AIAPILogsSourceType string

// AIAPILogsSourceType values
const (
	AIAPILogsSourceTypeChatbotDialog = AIAPILogsSourceType("chatbot_dialog") // ChatbotDialog
)
