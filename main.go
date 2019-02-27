package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/db"
	"github.com/vancelongwill/gotodos/handlers"
	"io/ioutil"
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
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret, jwtFound := env["JWT_SECRET"]

	if !jwtFound {
		log.Fatal("JWT_SECRET must be present in .env")
	}

	db, dbErr := db.Init()
	if dbErr != nil {
		log.Fatal("Error initialising database:\t", dbErr)
	}

	app := gin.Default()
	app.GET("/ping", ping)

	apiPrefix := "api"
	version := "v1"
	port := "8080"

	todoHandler := handlers.NewTodoHandler(db, jwtSecret)

	todoRouter := app.Group(path.Join(apiPrefix, version, "todos"))
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	// authorized.Use(AuthRequired())
	{
		todoRouter.GET("/", todoHandler.GetAll)
		todoRouter.POST("/", todoHandler.Create)
		todoRouter.GET("/:id", todoHandler.Get)
		todoRouter.PUT("/:id", todoHandler.Update)
		todoRouter.DELETE("/:id", todoHandler.Delete)
	}

	userHandler := handlers.NewUserHandler(db, jwtSecret)
	userRouter := app.Group(path.Join(apiPrefix, version, "user"))
	{
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/register", userHandler.Register)
	}

	app.Run(fmt.Sprintf(":%s", port))
}
