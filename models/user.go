package models

import (
	"database/sql"
)

// User defines the shape of an application user
type User struct {
	ID        uint
	FirstName sql.NullString
	LastName  sql.NullString
	Email     string
	Password  string
}

type UserStore interface {
	GetUser(email string) (*User, error)
	CreateUser(u *User) (*User, error)
}

// Serialize converts the user struct to a simple string map for conversion to JSON
func (u *User) Serialize() map[string]interface{} {
	mappedUser := map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
	}
	if u.FirstName.Valid {
		mappedUser["firstName"] = u.FirstName
	}
	if u.LastName.Valid {
		mappedUser["lastName"] = u.LastName
	}
	return mappedUser
}

// GetUser finds a single user from an sql database
func (db *DB) GetUser(email string) (*User, error) {
	sqlStatement := `SELECT * FROM users WHERE email = $1`
	u := new(User)
	err := db.QueryRow(sqlStatement, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// CreateUser inserts a single new user into an sql database
func (db *DB) CreateUser(u *User) (*User, error) {
	sqlStatement := `
	INSERT INTO users (email, first_name, last_name, password_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	err := db.QueryRow(sqlStatement, u.Email, u.FirstName, u.LastName, u.Password).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}
