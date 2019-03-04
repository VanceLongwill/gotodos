package models

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

// Todo defines the shape of a single todo item
type Todo struct {
	ID          uint
	Title       sql.NullString
	Note        sql.NullString
	CreatedAt   time.Time
	ModifiedAt  time.Time
	DueAt       pq.NullTime
	UserID      uint
	CompletedAt pq.NullTime
	IsDone      bool
}

const (
	resultsPerPage = 10
)

func (t *Todo) Serialize() map[string]interface{} {
	mappedTodo := map[string]interface{}{
		"id":         t.ID,
		"isDone":     t.IsDone,
		"createdAt":  t.CreatedAt,
		"modifiedAt": t.ModifiedAt,
	}

	if t.Title.Valid {
		mappedTodo["title"] = t.Title
	}
	if t.Note.Valid {
		mappedTodo["note"] = t.Note
	}
	if t.DueAt.Valid {
		mappedTodo["dueAt"] = t.DueAt
	}
	if t.CompletedAt.Valid {
		mappedTodo["completedAt"] = t.CompletedAt
	}

	return mappedTodo
}

func CreateTodo(db *sql.DB, t *Todo) error {
	sqlStatement := `
	INSERT INTO todos (title, note, created_at, modified_at, due_at, user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;`

	var newTodo Todo
	currentTime := time.Now()
	err := db.QueryRow(sqlStatement, t.Title, t.Note, currentTime, currentTime, currentTime, t.UserID).
		Scan(&newTodo.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTodos(db *sql.DB, userID, previousID uint) ([]*Todo, error) {
	sqlStatement := `
	SELECT * FROM todos WHERE user_id = $1 AND id > $3
	LIMIT $2;`
	rows, err := db.Query(sqlStatement, userID, resultsPerPage, previousID)

	defer rows.Close()

	todos := make([]*Todo, 0)
	for rows.Next() {
		t := new(Todo)
		if err := rows.
			Scan(&t.ID, &t.Title, &t.Note, &t.CreatedAt, &t.ModifiedAt, &t.DueAt, &t.UserID, &t.CompletedAt, &t.IsDone); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func GetTodo(db *sql.DB, todoID, userID uint) (*Todo, error) {
	sqlStatement := "SELECT * FROM todos WHERE id = $1 AND user_id = $2"
	todo := new(Todo)
	err := db.QueryRow(sqlStatement, todoID, userID).
		Scan(&todo.ID, &todo.Title, &todo.Note, &todo.CreatedAt,
			&todo.ModifiedAt, &todo.DueAt, &todo.UserID, &todo.CompletedAt, &todo.IsDone)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func MarkTodoAsComplete(db *sql.DB, todoID, userID uint) (*Todo, error) {
	sqlStatement := `
	UPDATE todos
	SET completed_at = $3, is_done = $4,
	WHERE id = $1 AND user_id = $2;`

	currentTime := time.Now()
	todo := new(Todo)

	execErr := db.QueryRow(sqlStatement, todoID, userID, todo.Title, todo.Note, currentTime).
		Scan(&todo.ID)
	if execErr != nil {
		return nil, execErr
	}
	return todo, nil
}

func UpdateTodo(db *sql.DB, t *Todo) (*Todo, error) {
	sqlStatement := `
	UPDATE todos
	SET title = $3, note = $4, modified_at = $5,
	WHERE id = $1 AND user_id = $2;`

	currentTime := time.Now()

	todo := new(Todo)

	execErr := db.QueryRow(sqlStatement, todo.ID, todo.UserID, todo.Title, todo.Note, currentTime).
		Scan(&todo.ID)
	if execErr != nil {
		return nil, execErr
	}
	return todo, nil
}

func DeleteTodo(db *sql.DB, todoID, userID uint) (uint, error) {
	sqlStatement := `
	DELETE FROM todos
	WHERE id = $1 AND user_id = $2;`
	_, err := db.Exec(sqlStatement, todoID, userID)
	if err != nil {
		return 0, err
	}
	return todoID, nil
}
