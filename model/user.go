package model

import (
	"time"
)

// User defines the shape of an application user
type User struct {
	ID        string    `json: "id"`
	Email     string    `json: "email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json: "createdAt"`
	DueAt     time.Time `json: "dueAt"`
}

// UserStorer defines a data store agnostic interface for CRUD operations
type UserStorer interface {
	GetAll() ([]User, error)
	Get(id User.ID) (User, error)
	Create(User) error
	Update(User) error
	Delete(id User.ID) error
}
