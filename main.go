package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/db"
	"github.com/vancelongwill/gotodos/handlers"
	"github.com/vancelongwill/gotodos/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const (
	apiVersion = "v1"
	apiPrefix  = "api"
	apiPort    = 8080
)

func getSecret() []byte {
	fn := "jwtsecret.key"
	key, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal("Error reading from: ", fn, err)
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
	db, dbErr := db.Init("0.0.0.0", "5432", "gotodos", "gotodos", "gotodos")
	if dbErr != nil {
		log.Fatal("Error initialising database:\t", dbErr)
	}

	jwtSecret := getSecret()
	app := gin.Default()
	app.GET("/ping", ping)

	todoHandler := handlers.NewTodoHandler(db, jwtSecret)

	todoRouter := app.Group(path.Join(apiPrefix, apiVersion, "todos"))

	todoRouter.Use(middleware.Authorize(jwtSecret))
	{
		todoRouter.GET("/", todoHandler.GetAll)
		todoRouter.POST("/", todoHandler.Create)
		todoRouter.GET("/:id", todoHandler.Get)
		todoRouter.PUT("/:id", todoHandler.Update)
		todoRouter.DELETE("/:id", todoHandler.Delete)
	}

	userHandler := handlers.NewUserHandler(db, jwtSecret)
	userRouter := app.Group(path.Join(apiPrefix, apiVersion, "user"))
	{
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/register", userHandler.Register)
	}

	app.Run(fmt.Sprintf(":%d", apiPort))
}
