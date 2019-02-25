package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/db"
	"github.com/vancelongwill/gotodos/handlers"
	"log"
	"net/http"
	"path"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {
	app := gin.Default()
	app.GET("/ping", ping)

	apiPrefix := "api"
	version := "v1"
	port := "8080"

	db, dbErr := db.Init()
	if dbErr != nil {
		log.Fatal("Error initialising database:\t", dbErr)
	}

	todoHandler := handlers.NewTodoHandler(db)

	todoRouter := app.Group(path.Join(apiPrefix, version, "todos"))
	{
		todoRouter.GET("/", todoHandler.GetAll)
		todoRouter.POST("/", todoHandler.Create)
		todoRouter.GET("/:id", todoHandler.Get)
		todoRouter.PUT("/:id", todoHandler.Update)
		todoRouter.DELETE("/:id", todoHandler.Delete)
	}

	app.Run(fmt.Sprintf(":%s", port))
}
