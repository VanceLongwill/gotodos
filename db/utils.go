package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver for gorm
	"github.com/vancelongwill/gotodos/models"
)

// Init starts the database
func Init() (*gorm.DB, error) {
	// open a db connection
	db, err := gorm.Open("postgres", "host=0.0.0.0 port=5432 user=gotodos dbname=gotodos password=gotodos sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&models.Todo{}, &models.User{})
	db.Model(&models.Todo{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.LogMode(true)

	return db, nil
}
