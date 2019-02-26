package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

// Todo defines the shape of a single todo item
type Todo struct {
	gorm.Model
	Title  string    `json: "title"`
	Note   string    `json: "note"`
	DueAt  time.Time `json: "dueAt"`
	IsDone bool      `json: "isDone"`
	UUID   string    `json: "id"`
}

func (todo *Todo) BeforeCreate(scope *gorm.Scope) error {
	uuid, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		return uuidErr
	}
	scope.SetColumn("UUID", uuid)
	return nil
}
