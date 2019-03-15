package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/models"
	"net/http"
	"strconv"
)

// TodoHandler wraps all handlers for Todos
type TodoHandler struct {
	db     *sql.DB
	secret []byte
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(db *sql.DB, secret []byte) *TodoHandler {
	return &TodoHandler{
		db:     db,
		secret: secret,
	}
}

func stringToUint(n string) uint {
	u64, err := strconv.ParseUint(n, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(u64)
}

type CreateBody struct {
	Title string `json: "title" binding: "required"`
	Note  string `json: "note" binding: "required"`
}

// Create creates a new item of type Todo and stores it
func (t *TodoHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var body CreateBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad request",
		})
		return
	}

	todo := models.Todo{
		Title:  sql.NullString{body.Title, true},
		Note:   sql.NullString{body.Note, true},
		UserID: userID,
		IsDone: false,
	}

	if err := models.CreateTodo(t.db, &todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     http.StatusCreated,
			"message":    "Unable to save todo",
			"resourceId": todo.ID,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Todo item created successfully!",
		"resourceId": todo.ID,
	})
}

// GetAll returns a all the current User's Todos
func (t *TodoHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	prev := c.DefaultQuery("prev", "0") // use previous id for pagination
	previousID := stringToUint(prev)

	todos, err := models.GetAllTodos(t.db, userID, previousID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Error fetching todos",
		})
		return
	}

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No todo items",
		})
		return
	}

	data := make([]map[string]interface{}, len(todos))
	for i, item := range todos {
		data[i] = item.Serialize()
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": data})
}

// Get returns a single Todo by ID
func (t *TodoHandler) Get(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	todoID := stringToUint(c.Param("id"))

	todo, err := models.GetTodo(t.db, todoID, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Unable to find todo",
		})
		return
	}

	if todo.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Todo doesn't belong to user",
		})
		return
	}

	_todo := todo.Serialize()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

type UpdateBody struct {
	Title string `json: "title"`
	Note  string `json: "note"`
}

// Update edits an existing Todo
func (t *TodoHandler) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	todoID := stringToUint(c.Param("id"))

	var body UpdateBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad request",
		})
		return
	}

	_todo := models.Todo{
		ID:     todoID,
		UserID: userID,
		Title:  sql.NullString{body.Title, true},
		Note:   sql.NullString{body.Note, true},
	}

	todo, err := models.UpdateTodo(t.db, _todo)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find todo", "resourceId": todo.ID})
		return
	}

	// completed, _ := strconv.Atoi(c.PostForm("completed"))
	// t.db.Model(&todo).Update("completed", false) // @TODO: add ability to change IsDone status
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// Delete deletes a single Todo
func (t *TodoHandler) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	todoID := stringToUint(c.Param("id"))

	deletedTodoID, err := models.DeleteTodo(t.db, todoID, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!", "resourceId": deletedTodoID})
}
