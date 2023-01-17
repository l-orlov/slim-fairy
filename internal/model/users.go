package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Age  int       `db:"age"`
}
