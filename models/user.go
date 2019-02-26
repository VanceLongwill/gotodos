package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// User defines the shape of an application user
type User struct {
	gorm.Model
	Email     string `json: "email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	UUID      string `json: "id"`
	Todos     []Todo // gorm will make a hasMany association based on User.ID & Todo.UserID (default behaviour)
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	uuid, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		return uuidErr
	}
	scope.SetColumn("UUID", uuid)
	return nil
}
