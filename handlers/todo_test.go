package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStringToUint(t *testing.T) {
	t.Log(`Should convert strings to uint values`)
	tests := []struct {
		input            string
		expected         uint
		shouldRaiseError bool
	}{
		{"156", 156, false},
		{"0", 0, false},
		{"1", 1, false},
		{"180180", 180180, false},
		{"-*", 0, true},
	}

	for _, test := range tests {
		res, err := StringToUint(test.input)
		if err != nil && !test.shouldRaiseError {
			t.Errorf("Unexpected error %s", err.Error())
			t.Fail()
		} else if err == nil && test.shouldRaiseError {
			t.Errorf("Expected error to be returned")
			t.Fail()
		}

		if !test.shouldRaiseError && res != test.expected {
			t.Errorf("Expected %d but received %d", test.expected, res)
		}
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)

	_, ok := getUserIDFromContext(mockContext)
	if ok {
		t.Errorf("should fail when userID is not set")
		t.Fail()
		return
	}

	mockContext.Set("userID", uint(1))
	userID, ok := getUserIDFromContext(mockContext)
	if !ok {
		t.Errorf("unable to get userID")
		t.Fail()
		return
	}

	if userID != 1 {
		t.Errorf("Expected ID to be 1, instead found: %d", userID)
		t.Fail()
	}
}

func TestGetTodoIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)

	_, ok := getTodoIDFromContext(mockContext)
	if ok {
		t.Errorf("should return false when todoID is not set")
		t.Fail()
		return
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d but received %d", http.StatusBadRequest, recorder.Code)
	}

	recorder = httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ = gin.CreateTestContext(recorder)
	mockContext.Params = []gin.Param{{Key: "id", Value: "1"}}

	todoID, ok := getTodoIDFromContext(mockContext)
	if !ok {
		t.Errorf("unable to get todoID")
		t.Fail()
		return
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d but received %d", http.StatusOK, recorder.Code)
	}

	if todoID != 1 {
		t.Errorf("Expected ID to be 1, instead found: %d", todoID)
		t.Fail()
	}

	recorder = httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ = gin.CreateTestContext(recorder)
	mockContext.Params = []gin.Param{{Key: "id", Value: "*&*&^-"}}
	_, ok = getTodoIDFromContext(mockContext)
	if ok {
		t.Errorf(`should return false when id is unconvertable`)
		t.Fail()
		return
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d but received %d", http.StatusBadRequest, recorder.Code)
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

func (db mockComplete) MarkTodoAsComplete(todoID, userID uint, currentTime time.Time) error {
	return nil
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
		t.Title = models.MakeNullString(fmt.Sprintf("TODO NO. %d", i))
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

func TestCreateTodo(t *testing.T) {
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

func TestGetTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)
	userID := uint(11)
	mockContext.Set("userID", userID)
	mockContext.Params = []gin.Param{{Key: "id", Value: "1"}}

	db := mockGet{}

	GetTodo(db)(mockContext)

	if recorder.Code != 200 {
		t.Errorf("Expected status code %d but received %d",
			200, recorder.Code)
		t.Log(recorder.Body)
		t.Fail()
	}
}

func TestUpdateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []mock{
		{
			`{"title": "test title","note": "test note"}`,
			http.StatusOK,
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
		mockContext.Set("userID", uint(1))
		mockContext.Params = []gin.Param{{Key: "id", Value: "1"}}
		req, _ := http.NewRequest("POST", "http://example.com/",
			bytes.NewBuffer([]byte(test.json)))
		mockContext.Request = req
		db := mockUpdate{}
		UpdateTodo(db)(mockContext)

		if recorder.Code != test.expectedCode {
			t.Errorf("Expected status code %d but received %d",
				test.expectedCode, recorder.Code)
			t.Log(recorder.Body)
		}
	}
}

func TestDeleteTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)
	mockContext.Set("userID", uint(1))
	mockContext.Params = []gin.Param{{Key: "id", Value: "1"}}
	db := mockDelete{}
	DeleteTodo(db)(mockContext)

	if recorder.Code != 200 {
		t.Errorf("Expected status code %d but received %d",
			200, recorder.Code)
		t.Log(recorder.Body)
	}
}

func TestMarkTodoAsComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder() // implements http.ResponseWriter
	mockContext, _ := gin.CreateTestContext(recorder)
	mockContext.Set("userID", uint(1))
	mockContext.Params = []gin.Param{{Key: "id", Value: "1"}}
	db := mockComplete{}
	MarkTodoAsComplete(db)(mockContext)

	if recorder.Code != 200 {
		t.Errorf("Expected status code %d but received %d",
			200, recorder.Code)
		t.Log(recorder.Body)
	}
}
