package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Init(host, port, user, password, dbname string) (*sql.DB, error) {
	// open a db connection
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}
	if connectErr := db.Ping(); connectErr != nil {
		return nil, connectErr
	}
	return db, nil
}
