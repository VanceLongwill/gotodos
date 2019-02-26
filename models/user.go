package models

import "github.com/jinzhu/gorm"

// User defines the shape of an application user
type User struct {
	gorm.Model
	Email     string `json: "email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Todos     []Todo // gorm will make a hasMany association based on User.ID & Todo.UserID (default behaviour)
}
