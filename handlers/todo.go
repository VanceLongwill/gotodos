package handlers

import (
	// "encoding/json"
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver for gorm
	"github.com/vancelongwill/gotodos/models"
	"net/http"
	"time"
	// "strconv"
)

// TodoHandler wraps all handlers for Todos
type TodoHandler struct {
	db     *gorm.DB
	secret string
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(db *gorm.DB, secret string) *TodoHandler {
	return &TodoHandler{
		db:     db,
		secret: secret,
	}
}

// Create creates a new item of type Todo and stores it
func (t *TodoHandler) Create(c *gin.Context) {
	uuid, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to save todo"})
		return
	}
	todo := models.Todo{
		Title:  c.PostForm("title"), // @TODO sanitize
		Note:   c.PostForm("note"),  // @TODO sanitize
		UUID:   uuid.String(),
		IsDone: false,
	}
	t.db.NewRecord(todo)
	t.db.Create(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item created successfully!",
		"resourceId": todo.UUID,
	})
}

type transformedTodo struct {
	ID     string    `json: "id"`
	Title  string    `json: "title"`
	Note   string    `json: "note"`
	DueAt  time.Time `json: "dueAt"`
	IsDone bool      `json: "isDone"`
}

// GetAll returns a all the current User's Todos
func (t *TodoHandler) GetAll(c *gin.Context) {
	var todos []models.Todo

	t.db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No todo items",
		})
		return
	}

	data := make([]transformedTodo, len(todos))
	for i, item := range todos {
		data[i] = transformedTodo{
			ID:     item.UUID,
			Title:  item.Title,
			Note:   item.Note,
			DueAt:  item.DueAt,
			IsDone: false, // @TODO: remove hardcoded IsDone
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": data})
}

// Get returns a single Todo by ID
func (t *TodoHandler) Get(c *gin.Context) {
	var todo models.Todo
	todoID := c.Param("id")

	t.db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Unable to find todo",
		})
		return
	}

	_todo := transformedTodo{
		ID:     todo.UUID,
		Title:  todo.Title,
		DueAt:  todo.DueAt,
		IsDone: false, // @TODO: remove hardcoded IsDone
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// Update edits an existing Todo
func (t *TodoHandler) Update(c *gin.Context) {
	var todo models.Todo
	todoID := c.Param("id")

	t.db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find todo"})
		return
	}

	t.db.Model(&todo).Update("title", c.PostForm("title"))
	t.db.Model(&todo).Update("note", c.PostForm("note"))
	// completed, _ := strconv.Atoi(c.PostForm("completed"))
	// t.db.Model(&todo).Update("completed", false) // @TODO: add ability to change IsDone status
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// Delete deletes a single Todo
func (t *TodoHandler) Delete(c *gin.Context) {
	var todo models.Todo
	todoID := c.Param("id")

	t.db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find todo"})
		return
	}

	t.db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
