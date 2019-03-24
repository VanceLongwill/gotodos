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
func NewDB(host, port, user, password, dbname string) (*DB, error) {
	// open a db connection
	connectstring := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectstring)
	if err != nil {
		return nil, err
	}
	if connecterr := db.Ping(); connecterr != nil {
		return nil, connecterr
	}
	return &DB{db}, nil
}
