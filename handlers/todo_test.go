package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStringToUint(t *testing.T) {
	t.Log(`Should convert strings to uint values`)
	tests := []struct {
		input    string
		expected uint
	}{
		{"156", 156},
		{"0", 0},
		{"1", 1},
		{"180180", 180180},
	}

	for _, test := range tests {
		res := StringToUint(test.input)
		if res != test.expected {
			t.Errorf("Expected %d but received %d", test.expected, res)
		}
	}
}

type mockDB struct{}

func (db mockDB) CreateTodo(t *Todo) error {
	return nil
}

// func GetAllTodos(userID, previousID uint) ([]*Todo, error)  {}
func (db mockDB) GetTodo(todoID, userID uint) (*Todo, error) {
	return &Todo{
		ID:     todoID,
		UserID: userID,
	}, nil
}

func (db mockDB) DeleteTodo(todoID, userID uint) (uint, error) {
	return todoID, nil
}
func (db mockDB) MarkTodoAsComplete(todoID, userID uint) (*Todo, error) {
	return &Todo{
		ID:     todoID,
		UserID: userID,
	}, nil
}
func (db mockDB) UpdateTodo(t Todo) (*Todo, error) {
	return &t, nil
}

func (db mockDB) GetAllTodos(id, prevID uint) ([]*Todo, error) {
	todos := make([]*Todo, 0)
	for i := uint(0); i < 10; i++ {
		t := &Todo{}
		t.ID = i
		t.Title = sql.NullString{fmt.Sprintf("TODO NO. %d", i), true}
		todos = append(todos, t)
	}
	return todos, nil
}
func TestGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)

	// Bug in gin.Context.DefaultQuery with CreateTestContext
	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	mockContext.Request = req

	db := &mockDB{}
	todoHandler := &TodoHandler{
		db:     db,
		secret: []byte("some-secret"),
	}

	userID := uint(11)
	mockContext.Set("userID", userID)

	// without userID in context
	todoHandler.GetAll(mockContext)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d but received %d", http.StatusOK, recorder.Code)
	}

	// res := recorder.Result()
	// res.Write
	// if err := res.Write(os.Stdout); err != nil {
	// 	t.Error(err)
	// }
}

func TestGetAllWithoutUserID(t *testing.T) {
	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)

	// Bug in gin.Context.DefaultQuery with CreateTestContext
	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	mockContext.Request = req

	db := mockDB{}
	todoHandler := &TodoHandler{
		db:     db,
		secret: []byte("some-secret"),
	}

	// without userID in context
	todoHandler.GetAll(mockContext)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d but received %d", http.StatusInternalServerError, recorder.Code)
	}

	// res := recorder.Result()
	// res.Write
	// if err := res.Write(os.Stdout); err != nil {
	// 	t.Error(err)
	// }
}

// func TestCreate(t *testing.T) {
// 	mockContext, _ := gin.CreateTestContext()
// 	todoHandler := NewTodoHandler(db, []byte("secret"))
// 	todoBody := &CreateBody{
// 		Title: "some title",
// 		Note:  "some note",
// 	}
//
// 	todoHandler.Create()
//
// }
