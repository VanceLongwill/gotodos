package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	// "github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"strings"
	"testing"
	"time"
)

type TodoTest struct {
	Todo     Todo
	Expected string
}

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

func TestCreateTodo(t *testing.T) {
	t.Log(`Should fail to create an empty todo`)
	// db, _, err := sqlmock.New()
	// if err != nil {
	// 	t.Error(err)
	// 	t.Fail()
	// }
	// defer db.Close()

	emptyTodo := Todo{
		Title: sql.NullString{"", true},
		Note:  sql.NullString{"", true},
	}

	if err := CreateTodo(db, &emptyTodo); err == nil {
		t.Error(err)
		t.Fail()
	}

	t.Log(`Should create a normal todo with no title`)
	mockNow := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	todo := Todo{
		ID:          1,
		Title:       sql.NullString{"", false},
		Note:        sql.NullString{"example note", true},
		CreatedAt:   mockNow,
		ModifiedAt:  mockNow,
		DueAt:       pq.NullTime{mockNow, false},
		UserID:      22,
		CompletedAt: pq.NullTime{mockNow, false},
		IsDone:      false,
	}

	if err := CreateTodo(db, &todo); err != nil {
		t.Error(err)
	}

}
func TestGetAllTodos(t *testing.T) {
	t.Log(`- GET ALL`)

}
func TestGetTodo(t *testing.T) {

}
func TestMarkTodoAsComplete(t *testing.T) {

}
func TestUpdateTodo(t *testing.T) {

}
func TestDeleteTodo(t *testing.T) {

}
