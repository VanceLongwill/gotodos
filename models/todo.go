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

const (
	timestampFormat = "2016-06-22 19:10:25"
)

func CreateTodo(db *sql.DB, t *Todo) (*Todo, error) {
	sqlStatement := `
	INSERT INTO todos (title, note, created_at, modified_at, due_at, user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`

	var newTodo Todo
	currentTime := time.Now()
	formattedTime := currentTime.Format(timestampFormat)
	err := db.QueryRow(sqlStatement, t.Title, t.Note, formattedTime, formattedTime, formattedTime).Scan(&newTodo.ID)
	if err != nil {
		return nil, err
	}
	return &newTodo, nil
}

func GetAllTodos(db *sql.DB, userID uint) ([]*Todo, error) {
	sqlStatement := `
	SELECT * FROM todos where user_id = $1
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
	err := db.QueryRow(sqlStatement, t.ID, t.UserID).Scan(&t.ID, &t.Title, &t.Note, &t.CreatedAt, &t.ModifiedAt, &t.DueAt, &t.UserID, &t.IsDone)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func UpdateTodo(db *sql.DB, t *Todo) (*Todo, error) {
	sqlStatement := `
	UPDATE todos
	SET title = $3, note = $4, modified_at = $5,
	WHERE id = $1 AND user_id = $2;`

	currentTime := time.Now()
	formattedTime := currentTime.Format(timestampFormat)

	execErr := db.QueryRow(sqlStatement, t.ID, t.UserID, t.Title, t.Note, formattedTime).Scan(&t.ID, &t.Title, &t.Note, &t.CreatedAt, &t.ModifiedAt, &t.DueAt, &t.UserID, &t.IsDone)
	if execErr != nil {
		return nil, execErr
	}
	return t, nil
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
