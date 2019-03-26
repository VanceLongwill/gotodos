package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/models"
	"net/http"
	"strconv"
)

// Context allows mocking of the gin.Context functions required
// type Context interface {
// 	Get(s string) (string, bool)
// 	BindJSON(o interface{}) error
// 	JSON(code int, obj interface{})
// 	DefaultQuery(key, defaultValue string) string
// 	Param(s string) string
// }

// Todo is a type alias for convenience
type Todo = models.Todo

// TodoStore is the datastore API for todos
type TodoStore interface {
	DBCreateTodo
	DBGetTodo
	DBGetAllTodos
	DBDeleteTodo
	DBMarkTodoAsComplete
	DBUpdateTodo
}

// StringToUint is a util function which converts strings to uints
func StringToUint(n string) uint {
	u64, err := strconv.ParseUint(n, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(u64)
}

// DBCreateTodo represents the part of the datalayer responsible for creation
type DBCreateTodo interface {
	CreateTodo(t *Todo) error
}

type CreateRequestBody struct {
	Title string `json:"title"`
	Note  string `json:"note" binding:"required"`
}

// CreateTodo returns a function which handles requests to create todos
func CreateTodo(db DBCreateTodo) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Unable to get UserID",
			})
			return
		}
		userID := strUserID.(uint)

		var body CreateRequestBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Bad request: %s", err.Error()),
			})
			return
		}

		todo := Todo{
			Title:  sql.NullString{body.Title, true},
			Note:   sql.NullString{body.Note, true},
			UserID: userID,
		}

		if err := db.CreateTodo(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":     http.StatusInternalServerError,
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
}

// DBGetAllTodos represents the part of the datalayer responsible for getting a list of todos
type DBGetAllTodos interface {
	GetAllTodos(userID, previousID uint) ([]*Todo, error)
}

// GetAllTodos returns a function responsible for handling requests for all the current User's todos
func GetAllTodos(db DBGetAllTodos) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Unable to get UserID",
			})
			return
		}
		userID := strUserID.(uint)
		prev := c.DefaultQuery("prev", "0") // use previous id for pagination
		previousID := StringToUint(prev)

		todos, err := db.GetAllTodos(userID, previousID)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
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
}

// DBGetTodo represents the part of the datalayer responsible for getting a single todo
type DBGetTodo interface {
	GetTodo(todoID, userID uint) (*Todo, error)
}

// GetTodo returns a function which handles requests to get a single todo
func GetTodo(db DBGetTodo) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Unable to get UserID",
			})
			return
		}
		userID := strUserID.(uint)
		strTodoID := c.Param("id")
		if len(strTodoID) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Unable to get todo ID",
			})
			return
		}
		todoID := StringToUint(strTodoID)

		todo, err := db.GetTodo(todoID, userID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Unable to find todo",
			})
			return
		}

		_todo := todo.Serialize()
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
	}
}

// DBUpdateTodo represents the part of the datalayer responsible for updating a todo
type DBUpdateTodo interface {
	UpdateTodo(t Todo) (*Todo, error)
}

type UpdateRequestBody struct {
	Title string `json:"title"`
	Note  string `json:"note"`
}

// UpdateTodo returns a function which handles requests to edit an existing Todo
func UpdateTodo(db DBUpdateTodo) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Unable to get UserID",
			})
			return
		}
		userID := strUserID.(uint)
		strTodoID := c.Param("id")
		if len(strTodoID) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Unable to get todo ID",
			})
			return
		}
		todoID := StringToUint(strTodoID)

		var body UpdateRequestBody

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Bad request",
			})
			return
		}

		if len(body.Title) == 0 && len(body.Note) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Todo title or note must be provided to update",
			})
			return
		}

		_todo := Todo{
			ID:     todoID,
			UserID: userID,
			Title:  sql.NullString{body.Title, true},
			Note:   sql.NullString{body.Note, true},
		}

		todo, err := db.UpdateTodo(_todo)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find todo", "resourceId": todo.ID})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
	}
}

// DBDeleteTodo represents the part of the datalayer responsible for deleting a single todo
type DBDeleteTodo interface {
	DeleteTodo(todoID, userID uint) (uint, error)
}

// DeleteTodo returns a function which handles requests to delete a todo
func DeleteTodo(db DBDeleteTodo) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Unable to get UserID",
			})
			return
		}
		userID := strUserID.(uint)
		strTodoID := c.Param("id")
		if len(strTodoID) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Unable to get todo ID",
			})
			return
		}
		todoID := StringToUint(strTodoID)

		deletedTodoID, err := db.DeleteTodo(todoID, userID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find todo"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!", "resourceId": deletedTodoID})
	}
}

// DBMarkTodoAsComplete represents the part of the datalayer responsible for updating the completion status of a todo
type DBMarkTodoAsComplete interface {
	MarkTodoAsComplete(todoID, userID uint) error
}

// MarkTodoAsComplete returns a function which handles requests to mark a todo as complete
func MarkTodoAsComplete(db DBMarkTodoAsComplete) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Unable to get UserID",
			})
			return
		}
		userID := strUserID.(uint)
		strTodoID := c.Param("id")
		if len(strTodoID) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Unable to get todo ID",
			})
			return
		}
		todoID := StringToUint(strTodoID)

		if err := db.MarkTodoAsComplete(todoID, userID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Unable to find todo",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo marked as complete successfully"})
	}
}
