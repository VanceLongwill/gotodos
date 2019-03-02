package models

import (
	"database/sql"
)

// User defines the shape of an application user
type User struct {
	ID        uint   `json: "id"` // `gorm:"type:bigint(20) unsigned auto_increment;primary_key"`
	Email     string `json: "email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	// Todos     []Todo // `gorm:"ForeignKey:UserID"`
}

func (u *User) Serialize() map[string]interface{} {
	mappedUser := map[string]interface{}{
		"id":        u.ID,
		"email":     u.Email,
		"firstName": u.FirstName,
		"LastName":  u.LastName,
	}
	return mappedUser
}

func GetUser(db *sql.DB, u *User) (*User, error) {
	sqlStatement := `
	SELECT * FROM users WHERE email = $1
	`
	err := db.QueryRow(sqlStatement, u.Email).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUser(db *sql.DB, u *User) (*User, error) {
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
