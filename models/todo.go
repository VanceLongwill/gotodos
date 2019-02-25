package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Todo defines the shape of a single todo item
type Todo struct {
	gorm.Model
	Title  string    `json: "title"`
	Body   string    `json: "body"`
	DueAt  time.Time `json: "dueAt"`
	IsDone bool      `json: "isDone"`
	UserID uint
}
