package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // postgres driver
)

// DB is the database
type DB struct {
	*sql.DB
}

// NewDB makes & tests a connection with the DB specified then returns it
func NewDB(user, password, dbname, host string) (*DB, error) {
	connectString := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		user, password, host, dbname)
	// open a db connection
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
