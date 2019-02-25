package model

import (
	"time"
)

// Todo defines the shape of a single todo item
type Todo struct {
	ID        string    `json: "id"`
	Title     string    `json: "title"`
	Note      string    `json: "note"`
	CreatedAt time.Time `json: "createdAt"`
	DueAt     time.Time `json: "dueAt"`
}

// TodoStorer defines a data store agnostic interface for CRUD operations
type TodoStorer interface {
	GetAll() ([]Todo, error)
	Get(id Todo.ID) (Todo, error)
	Create(Todo) error
	Update(Todo) error
	Delete(id Todo.ID) error
}
