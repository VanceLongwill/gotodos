package handlers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver for gorm
	"github.com/vancelongwill/gotodos/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	// "strconv"
)

func createToken(data map[string]interface{}, expireTime time.Time, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data":      data,
		"expiresAt": expireTime.Unix(),
	})
	return token.SignedString(secret)
}

// UserHandler wraps all handlers for application Users
type UserHandler struct {
	db     *gorm.DB
	secret string
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(db *gorm.DB, secret string) *UserHandler {
	return &UserHandler{
		db:     db,
		secret: secret,
	}
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

	u.db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User deleted successfully!"})
}

// @TODO: type all requests and responses like below, adding required fields etc & error handling for bad requests

// RegisterRequest specifies the request body shape for registering a new application user
type RegisterRequest struct {
	Email     string `json: "email" binding: "required"`
	Password  string `json: "password" binding: "required"`
	FirstName string `json: "firstName" binding: "required"`
	LastName  string `json: "lastName" binding: "required"`
}

// Register creates a new application user
func (u *UserHandler) Register(c *gin.Context) {
	var body RegisterRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad request",
		})
		return
	}

	var existingUser models.User
	if err := u.db.Where("email = ?", body.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to login"})
		return
	}

	hashBytes, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to register user"})
		return
	}

	user := models.User{
		Email:     body.Email,
		Password:  string(hashBytes),
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}

	u.db.NewRecord(user)
	u.db.Create(&user)

	serializedUser := map[string]interface{}{
		"id":        user.UUID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	}
	expiry := time.Now().Add(time.Hour * 24 * 7)

	token, tokenErr := createToken(serializedUser, expiry, u.secret)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to register user"})
		return
	}

	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"user":       serializedUser,
		"token":      token,
		"resourceId": user.UUID,
		"message":    "User registered successfully!",
	})
}

// LoginRequest specifies the request body shape for logging in an application user
type LoginRequest struct {
	Email    string `json: "email" binding: "required"`
	Password string `json: "password" binding: "required"`
}

// Login provides the client with a jwt token to authenticate on other routes
func (u *UserHandler) Login(c *gin.Context) {
	var body LoginRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad request",
		})
		return
	}

	var user models.User
	if err := u.db.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to login"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unable to login"})
		return
	}

	serializedUser := map[string]interface{}{
		"id":        user.UUID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	}

	expiry := time.Now().Add(time.Hour * 24 * 7)
	token, tokenErr := createToken(serializedUser, expiry, u.secret)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to login"})
		return
	}

	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"user":    serializedUser,
		"token":   token,
		"message": "User logged in successfully!",
	})
}
