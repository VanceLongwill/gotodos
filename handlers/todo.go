package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver for gorm
	"github.com/vancelongwill/gotodos/models"
	"net/http"
	"time"
	// "strconv"
)

// TodoHandler wraps all handlers for Todos
type TodoHandler struct {
	db *gorm.DB
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db}
}

// Create creates a new item of type Todo and stores it
func (t *TodoHandler) Create(c *gin.Context) {
	todo := models.Todo{
		Title:  c.PostForm("title"), // @TODO sanitize
		Body:   c.PostForm("body"),  // @TODO sanitize
		IsDone: false,
	}
	t.db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item created successfully!",
		"resourceId": todo.ID,
	})
}

type transformedTodo struct {
	ID     uint      `json: "id"`
	Title  string    `json: "title"`
	Body   string    `json: "body"`
	DueAt  time.Time `json: "dueAt"`
	IsDone bool      `json: "isDone"`
}

// GetAll returns a all the current User's Todos
func (t *TodoHandler) GetAll(c *gin.Context) {
	var todos []models.Todo
	var _todos []transformedTodo

	t.db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Unable to find todo",
		})
		return
	}

	for _, item := range todos {
		_todos = append(_todos, transformedTodo{
			ID:     item.ID,
			Title:  item.Title,
			Body:   item.Body,
			DueAt:  item.DueAt,
			IsDone: false, // @TODO: remove hardcoded IsDone
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
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
		ID:     todo.ID,
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
	t.db.Model(&todo).Update("body", c.PostForm("body"))
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