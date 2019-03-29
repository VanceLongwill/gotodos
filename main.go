package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/handlers"
	"github.com/vancelongwill/gotodos/middleware"
	"github.com/vancelongwill/gotodos/models"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
)

// Env defines the environment variables necessary for the app to run
type Env struct {
	APIVersion       string `env:"API_VERSION"`
	APIPort          string `env:"API_PORT"`
	JWTSecret        string `env:"JWT_SECRET"`
	APIMode          string `env:"API_MODE"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresName     string `env:"POSTGRES_NAME"`
	PostgresHost     string `env:"POSTGRES_HOST"`
}

// getEnv gets all the necessary environment variables
func getEnv() *Env {
	env := Env{}
	v := reflect.ValueOf(&env).Elem()
	t := reflect.TypeOf(env)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("env")
		val := os.Getenv(tag)
		if len(val) == 0 {
			panic(fmt.Sprintf("Error getting environment variables: %s not set", tag))
		}
		v.Field(i).SetString(val)
	}
	return &env
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {
	// Load env vars into a struct
	env := getEnv()
	// Open DB connection
	db, dbErr := models.NewDB(env.PostgresUser,
		env.PostgresPassword, env.PostgresName, env.PostgresHost)
	if dbErr != nil {
		log.Fatal("Error initialising database:\t", dbErr)
	}

	app := gin.Default()
	app.GET("/ping", ping)

	jwtSecret := []byte(env.JWTSecret)

	// todo resources
	todoRouter := app.Group(path.Join("api", env.APIVersion, "todos"))
	todoRouter.Use(middleware.Authorize(jwtSecret))
	{
		todoRouter.GET("/", handlers.GetAllTodos(db))
		todoRouter.POST("/", handlers.CreateTodo(db))
		todoRouter.GET("/:id", handlers.GetTodo(db))
		todoRouter.PUT("/:id", handlers.UpdateTodo(db))
		todoRouter.GET("/:id/completed", handlers.MarkTodoAsComplete(db))
		todoRouter.DELETE("/:id", handlers.DeleteTodo(db))
	}

	// user resources
	userRouter := app.Group(path.Join("api", env.APIVersion, "user"))
	{
		userRouter.POST("/login", handlers.LoginUser(db, jwtSecret))
		userRouter.POST("/register", handlers.RegisterUser(db, jwtSecret))
	}

	app.Run(fmt.Sprintf(":%s", env.APIPort))
}
