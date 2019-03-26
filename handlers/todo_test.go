package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

type mockCreate struct{}

func (db mockCreate) CreateTodo(t *Todo) error {
	return nil
}

type mockGet struct{}

func (db mockGet) GetTodo(todoID, userID uint) (*Todo, error) {
	return &Todo{
		ID:     todoID,
		UserID: userID,
	}, nil
}

type mockDelete struct{}

func (db mockDelete) DeleteTodo(todoID, userID uint) (uint, error) {
	return todoID, nil
}

type mockComplete struct{}

func (db mockComplete) MarkTodoAsComplete(todoID, userID uint, currentTime time.Time) (*Todo, error) {
	return &Todo{
		ID:          todoID,
		UserID:      userID,
		CompletedAt: pq.NullTime{currentTime, true},
	}, nil
}

type mockUpdate struct{}

func (db mockUpdate) UpdateTodo(t Todo) (*Todo, error) {
	return &t, nil
}

type mockGetAll struct{}

func (db mockGetAll) GetAllTodos(id, prevID uint) ([]*Todo, error) {
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

	db := &mockGetAll{}
	userID := uint(11)
	mockContext.Set("userID", userID)

	// without userID in context
	GetAllTodos(db)(mockContext)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d but received %d", http.StatusOK, recorder.Code)
	}
}

func TestGetAllWithoutUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)

	// Bug in gin.Context.DefaultQuery with CreateTestContext
	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	mockContext.Request = req

	db := mockGetAll{}

	// without userID in context
	GetAllTodos(db)(mockContext)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d but received %d",
			http.StatusInternalServerError, recorder.Code)
	}
}

type mock struct {
	json         string
	expectedCode int
}

func TestCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []mock{
		{
			`{"title": "test title","note": "test note"}`,
			http.StatusCreated,
		},
		{
			`{}`,
			http.StatusBadRequest,
		},
		{
			`{"some": "field"}`,
			http.StatusBadRequest,
		},
		{
			`{"title": "", "note": ""}`,
			http.StatusBadRequest,
		},
		{
			`{"title": "asd", "note": ""}`,
			http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		recorder := httptest.NewRecorder() // implements http.ResponseWriter
		mockContext, _ := gin.CreateTestContext(recorder)
		userID := uint(11)
		mockContext.Set("userID", userID)
		req, _ := http.NewRequest("POST", "http://example.com/",
			bytes.NewBuffer([]byte(test.json)))
		mockContext.Request = req
		db := mockCreate{}
		CreateTodo(db)(mockContext)

		if recorder.Code != test.expectedCode {
			t.Errorf("Expected status code %d but received %d",
				test.expectedCode, recorder.Code)
			t.Log(recorder.Body)
		}
	}
}

// func TestUpdateTodo(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	tests := []mock{
// 		{
// 			`{"title": "test title","note": "test note"}`,
// 			http.StatusCreated,
// 		},
// 		{
// 			`{}`,
// 			http.StatusBadRequest,
// 		},
// 		{
// 			`{"some": "field"}`,
// 			http.StatusBadRequest,
// 		},
// 		{
// 			`{"title": "", "note": ""}`,
// 			http.StatusBadRequest,
// 		},
// 		{
// 			`{"title": "asd", "note": ""}`,
// 			http.StatusBadRequest,
// 		},
// 	}
//
// 	for _, test := range tests {
// 		recorder := httptest.NewRecorder() // implements http.ResponseWriter
// 		mockContext, _ := gin.CreateTestContext(recorder)
// 		userID := uint(11)
// 		mockContext.Set("userID", userID)
// 		req, _ := http.NewRequest("PUT", "localhost:8080/api/v1/todos/123",
// 			bytes.NewBuffer([]byte(test.json)))
// 		mockContext.Request = req
// 		db := mockUpdate{}
// 		UpdateTodo(db)(mockContext)
//
// 		if recorder.Code != test.expectedCode {
// 			t.Errorf("Expected status code %d but received %d",
// 				test.expectedCode, recorder.Code)
// 			t.Log(recorder.Body)
// 		}
// 	}
// }
