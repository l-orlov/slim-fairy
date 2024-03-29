package store

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/l-orlov/slim-fairy/bot/internal/model"
)

// Column names for models
var (
	asteriskUsers,
	asteriskNutritionists,
	asteriskAuthData,
	asteriskChatBotDialogs,
	asteriskAIAPILogs string
)

func init() {
	// Init column names for models
	asteriskUsers = asterisk(model.User{})
	asteriskNutritionists = asterisk(model.Nutritionist{})
	asteriskAuthData = asterisk(model.AuthData{})
	asteriskChatBotDialogs = asterisk(model.ChatBotDialog{})
	asteriskAIAPILogs = asterisk(model.AIAPILog{})
}

type tableNameGetter interface {
	DbTable() string
}

// asterisk replace * in queries select(*) by column names (only for models without nesting)
func asterisk(a tableNameGetter) string {
	modelType := reflect.TypeOf(a)
	var columns []string
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		columnName, ok := field.Tag.Lookup("db")
		if !ok || columnName == "-" {
			continue
		}
		columns = append(columns, fmt.Sprintf("%s.%s", a.DbTable(), columnName))
	}
	return strings.Join(columns, ", ")
}
