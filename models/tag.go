package models

// @TODO: implement tagging as a method to organise user todos

//
// import (
// 	"time"
// )
//
// // Tag defines the shape of a named tag items, which can have many chidren Todos
// type Tag struct {
// 	ID        string    `json: "id"`
// 	Name      string    `json: "name"`
// 	CreatedAt time.Time `json: "createdAt"`
// 	Todos     []Todo.ID `json: "todos"`
// }
//
// // TagStorer defines a data store agnostic interface for CRUD operations
// type TagStorer interface {
// 	GetAll() ([]Tag, error)
// 	Get(id Tag.ID) (Tag, error)
// 	Create(Tag) error
// 	Update(Tag) error
// 	Delete(id Tag.ID) error
// }
