package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func createTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "create",
	})
}

func getAllTodos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "getAll",
	})
}

func getTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get",
	})
}

func updateTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "update",
	})
}

func deleteTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete",
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", ping)

	v1prefix := "/api/v1"

	todoRouter := r.Group(v1prefix + "/todos")
	{
		todoRouter.POST("/", createTodo)
		todoRouter.GET("/", getAllTodos)
		todoRouter.GET("/:id", getTodo)
		todoRouter.PUT("/:id", updateTodo)
		todoRouter.DELETE("/:id", deleteTodo)
	}

	r.Run()
}
