package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"strings"
	"testing"
	"time"
)

type TodoTest struct {
	Todo     Todo
	Expected string
}

var todoTableRows = []string{"id", "title", "note",
	"created_at", "modified_at", "due_at",
	"user_id", "completed_at", "is_done"}

func TestSerialize(t *testing.T) {
	t.Log(`Should serialize todos correctly`)
	mockNow := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	todo := Todo{
		ID:          1,
		Title:       sql.NullString{"example title", true},
		Note:        sql.NullString{"example note", true},
		CreatedAt:   mockNow,
		ModifiedAt:  mockNow,
		DueAt:       pq.NullTime{mockNow, false},
		UserID:      22,
		CompletedAt: pq.NullTime{mockNow, false},
		IsDone:      false,
	}

	todoTests := make([]TodoTest, 4)

	todoTests[0].Todo = todo
	todoTests[0].Expected = `{"createdAt":"2009-11-10T23:00:00Z","id":1,"isDone":false,"modifiedAt":"2009-11-10T23:00:00Z","note":"example note","title":"example title"}`

	todo.DueAt.Valid = true
	todoTests[1].Todo = todo
	todoTests[1].Expected = `{"createdAt":"2009-11-10T23:00:00Z","dueAt":"2009-11-10T23:00:00Z","id":1,"isDone":false,"modifiedAt":"2009-11-10T23:00:00Z","note":"example note","title":"example title"}`
	todo.DueAt.Valid = false

	todo.CompletedAt.Valid = true
	todo.IsDone = true
	todoTests[2].Todo = todo
	todoTests[2].Expected = `{"completedAt":"2009-11-10T23:00:00Z","createdAt":"2009-11-10T23:00:00Z","id":1,"isDone":true,"modifiedAt":"2009-11-10T23:00:00Z","note":"example note","title":"example title"}`
	todo.CompletedAt.Valid = false
	todo.IsDone = false

	todo.Title.Valid = false
	todoTests[3].Todo = todo
	todoTests[3].Expected = `{"createdAt":"2009-11-10T23:00:00Z","id":1,"isDone":false,"modifiedAt":"2009-11-10T23:00:00Z","note":"example note"}`
	todo.Title.Valid = true

	for _, test := range todoTests {
		serialTodo := test.Todo.Serialize()

		data, err := json.Marshal(serialTodo)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		received := string(data)

		if strings.Compare(test.Expected, received) != 0 {
			fmt.Println(" === EXPECTED ===\n", test.Expected)
			fmt.Println(" === RECEIVED ===\n", received)
			t.Fail()
		}
	}

}

func TestCreateTodoEmpty(t *testing.T) {
	t.Log(`Should fail to create an empty todo`)

	var mockDB *sql.DB
	db := DB{mockDB}

	emptyTodo := Todo{
		Title: sql.NullString{"", true},
		Note:  sql.NullString{"", true},
	}

	if err := db.CreateTodo(&emptyTodo); err != ErrorEmptyTodo {
		t.Error(err)
		t.Fail()
	}
}

func TestCreateTodo(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	todo := Todo{
		ID:     1,
		Title:  sql.NullString{"asd", true},
		Note:   sql.NullString{"example note", true},
		UserID: 22,
	}

	mock.ExpectExec(`INSERT INTO todos`).
		WithArgs(todo.Title, todo.Note, todo.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := DB{mockDB}

	t.Log(`Should create a normal todo with no title`)
	// mockNow := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	if err := db.CreateTodo(&todo); err != nil {
		t.Error(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllTodosEmpty(t *testing.T) {
	t.Log(`Should get an empty list of todos`)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	userID := uint(1)
	prevID := uint(0)

	rows := sqlmock.NewRows(todoTableRows)

	mock.ExpectQuery(`SELECT \* FROM todos WHERE.+`).
		WithArgs(userID, resultsPerPage, prevID).
		WillReturnRows(rows)

	db := DB{mockDB}

	todos, err := db.GetAllTodos(userID, prevID)
	if err != nil {
		t.Errorf("failed to get todo list")
		t.Fail()
	}
	if len(todos) != 0 {
		t.Errorf("failed to get empty todo list")
		t.Fail()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllTodos(t *testing.T) {
	t.Log(`Should get a list of todos`)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	userID := uint(1)
	prevID := uint(0)

	timestmp := time.Now()

	numberOfRows := 8

	rows := sqlmock.NewRows(todoTableRows)
	for i := 1; i < numberOfRows+1; i++ {
		rows.
			AddRow(i, "title 1", "hello", timestmp, timestmp, timestmp, userID, timestmp, false)
	}

	mock.ExpectQuery(`SELECT \* FROM todos WHERE.+`).
		WithArgs(userID, resultsPerPage, prevID).
		WillReturnRows(rows)

	db := DB{mockDB}

	todos, err := db.GetAllTodos(userID, prevID)
	if err != nil {
		t.Errorf("failed to get todo list")
		t.Fail()
	}
	if len(todos) != numberOfRows {
		t.Errorf("failed to get the correct number of todos")
		t.Fail()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTodo(t *testing.T) {
	t.Log(`Should get a single todos`)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	userID := uint(1)
	todoID := uint(1)

	timestmp := time.Now()

	rows := sqlmock.NewRows(todoTableRows).
		AddRow(1, "title 1", "hello", timestmp, timestmp, timestmp, userID, timestmp, false)

	mock.ExpectQuery(`SELECT \* FROM todos WHERE.+`).
		WithArgs(todoID, userID).
		WillReturnRows(rows)

	db := DB{mockDB}

	todo, err := db.GetTodo(todoID, userID)

	t.Log(todo.Serialize())

	if err != nil {
		t.Errorf("failed to get todo list")
		t.Fail()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMarkTodoAsComplete(t *testing.T) {
	t.Log(`Should mark a single todo as done`)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	userID := uint(1)
	todoID := uint(1)

	currentTime := time.Now()

	mock.ExpectExec(`UPDATE todos.+`).
		WithArgs(todoID, userID, currentTime).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := DB{mockDB}

	if err := db.MarkTodoAsComplete(todoID, userID, currentTime); err != nil {
		t.Errorf("Failed to mark todo as done: %s", err.Error())
		t.Fail()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateTodo(t *testing.T) {
	t.Log(`Should update a single todo's title and note`)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	todo := Todo{
		Title:      sql.NullString{"modified title", true},
		Note:       sql.NullString{"changed note", true},
		UserID:     1,
		ID:         1,
		ModifiedAt: time.Now(),
	}

	mock.ExpectExec(`UPDATE todos.+`).
		WithArgs(todo.ID, todo.UserID, todo.Title, todo.Note, todo.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := DB{mockDB}
	if _, err := db.UpdateTodo(todo); err != nil {
		t.Errorf("Failed to mark todo as done: %s", err.Error())
		t.Fail()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTodo(t *testing.T) {
	t.Log(`Should delete a single todo`)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()
	todoID := uint(1)
	userID := uint(1)

	mock.ExpectExec(`DELETE FROM todos.+`).
		WithArgs(todoID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := DB{mockDB}
	if _, err := db.DeleteTodo(todoID, userID); err != nil {
		t.Errorf("Failed to delete todo: %s", err.Error())
		t.Fail()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
