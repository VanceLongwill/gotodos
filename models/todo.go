package models

import (
	"database/sql"
	"time"
)

// Todo defines the shape of a single todo item
type Todo struct {
	ID         uint // `gorm:"type:bigint(20) unsigned auto_increment;primary_key"`
	Title      string
	Note       string
	CreatedAt  string
	ModifiedAt string
	DueAt      string
	UserID     uint
	IsDone     bool
}

func (t *Todo) Serialize() map[string]interface{} {
	mappedTodo := map[string]interface{}{
		"id":     t.ID,
		"title":  t.Title,
		"note":   t.Note,
		"dueAt":  t.DueAt,
		"isDone": false, // @TODO remove constant
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
	err := db.QueryRow(sqlStatement, t.Title, t.Note, currentTime, currentTime, currentTime, t.UserID).Scan(&newTodo.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTodos(db *sql.DB, userID uint) ([]*Todo, error) {
	sqlStatement := `
	SELECT * FROM todos WHERE user_id = $1
	LIMIT $2;`
	rows, err := db.Query(sqlStatement, userID, 10) //@TODO pagination

	defer rows.Close()

	todos := make([]*Todo, 0)
	for rows.Next() {
		t := new(Todo)
		if err := rows.Scan(&t.ID, &t.Title, &t.Note, &t.CreatedAt, &t.ModifiedAt, &t.DueAt, &t.UserID, &t.IsDone); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func GetTodo(db *sql.DB, t *Todo) (*Todo, error) {
	sqlStatement := "SELECT * FROM todos WHERE id = $1 AND user_id = $2"
	todo := new(Todo)
	err := db.QueryRow(sqlStatement, todo.ID, todo.UserID).Scan(&todo.ID, &todo.Title, &todo.Note, &todo.CreatedAt, &todo.ModifiedAt, &todo.DueAt, &todo.UserID, &todo.IsDone)
	if err != nil {
		return nil, err
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

	execErr := db.QueryRow(sqlStatement, todo.ID, todo.UserID, todo.Title, todo.Note, currentTime).Scan(&todo.ID, &todo.Title, &todo.Note, &todo.CreatedAt, &todo.ModifiedAt, &todo.DueAt, &todo.UserID, &todo.IsDone)
	if execErr != nil {
		return nil, execErr
	}
	return todo, nil
}

func DeleteTodo(db *sql.DB, t *Todo) (*Todo, error) {
	sqlStatement := `
	DELETE FROM todos
	WHERE id = $1 AND user_id = $2;`
	_, err := db.Exec(sqlStatement, t.ID, t.UserID)
	if err != nil {
		return nil, err
	}
	return t, nil
}
