package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// User defines the shape of an application user
type User struct {
	gorm.Model
	Email    string    `json: "email"`
	Password string    `json:"password"`
	DueAt    time.Time `json: "dueAt"`
	Todos    []Todo    // gorm will make a hasMany association based on User.ID & Todo.UserID (default behaviour)
}
