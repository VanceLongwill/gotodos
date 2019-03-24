package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/handlers"
	"github.com/vancelongwill/gotodos/middleware"
	"github.com/vancelongwill/gotodos/models"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const (
	apiVersion  = "v1"
	apiPrefix   = "api"
	apiPort     = 8080
	jwtSecretFn = "jwtsecret.key"
)

func getSecret() []byte {
	key, err := ioutil.ReadFile(jwtSecretFn)
	if err != nil {
		log.Fatal("Error reading from: ", jwtSecretFn, err)
	}
	if len(key) == 0 {
		log.Fatal("Empty jwt key")
	}
	return key
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {
	db, dbErr := models.NewDB("0.0.0.0", "5432", "gotodos", "gotodos", "gotodos")
	if dbErr != nil {
		log.Fatal("Error initialising database:\t", dbErr)
	}

	jwtSecret := getSecret()
	app := gin.Default()
	app.GET("/ping", ping)

	todoRouter := app.Group(path.Join(apiPrefix, apiVersion, "todos"))

	todoRouter.Use(middleware.Authorize(jwtSecret))
	{
		todoRouter.GET("/", handlers.GetAllTodos(db))
		todoRouter.POST("/", handlers.CreateTodo(db))
		todoRouter.GET("/:id", handlers.GetTodo(db))
		todoRouter.PUT("/:id", handlers.UpdateTodo(db))
		todoRouter.GET("/:id/completed", handlers.MarkTodoAsComplete(db))
		todoRouter.DELETE("/:id", handlers.DeleteTodo(db))
	}

	userRouter := app.Group(path.Join(apiPrefix, apiVersion, "user"))
	{
		userRouter.POST("/login", handlers.LoginUser(db, jwtSecret))
		userRouter.POST("/register", handlers.RegisterUser(db, jwtSecret))
	}

	app.Run(fmt.Sprintf(":%d", apiPort))
}
