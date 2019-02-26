package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver for gorm
	"github.com/vancelongwill/gotodos/models"
	"net/http"
	"time"
	// "strconv"
)

// UserHandler wraps all handlers for application Users
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db}
}

// Create creates a new application user
func (u *UserHandler) Create(c *gin.Context) {
	user := models.User{
		Email:     c.PostForm("email"),
		FirstName: c.PostForm("firstName"),
		LastName:  c.PostForm("lastName"),
		Password:  c.PostForm("password"),
	}
	u.db.Save(&user)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "User registered successfully!",
		"resourceId": user.ID,
	})
}

// Delete deletes a single application user
func (u *UserHandler) Delete(c *gin.Context) {
	var user models.User
	userID := c.Param("id")

	u.db.First(&user, userID)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find user"})
		return
	}

	t.db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User deleted successfully!"})
}
