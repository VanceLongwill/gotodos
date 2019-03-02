package handlers

import (
	"database/sql"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vancelongwill/gotodos/middleware"
	"github.com/vancelongwill/gotodos/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	// "strconv"
)

func createToken(userID uint, expireTime time.Time, secret string) (string, error) {
	type Claims = middleware.UserClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	})
	return token.SignedString([]byte(secret))
}

// UserHandler wraps all handlers for application Users
type UserHandler struct {
	db     *sql.DB
	secret string
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(db *sql.DB, secret string) *UserHandler {
	return &UserHandler{
		db:     db,
		secret: secret,
	}
}

// Delete deletes a single application user
// func (u *UserHandler) Delete(c *gin.Context) {
// 	var user models.User
// 	userID := c.Param("id")
//
// 	u.db.First(&user, userID)
//
// 	if user.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Unable to find user"})
// 		return
// 	}
//
// 	u.db.Delete(&user)
// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User deleted successfully!"})
// }

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

	// var existingUser models.User
	// if err := u.db.Where("email = ?", body.Email).First(&existingUser).Error; err == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Email already used"})
	// 	return
	// }

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

	newUser, err := models.CreateUser(u.db, &user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to register user"})
		return
	}
	expiryInterval := 60 * 60 * 24 * 7
	expiry := time.Now().Add(time.Hour * 24 * 7)

	token, tokenErr := createToken(user.ID, expiry, u.secret)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to register user"})
		return
	}

	c.SetCookie("token", token, expiryInterval, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"email":      user.Email,
		"token":      token,
		"resourceId": newUser.ID,
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
	reqUser := models.User{Email: body.Email, Password: body.Password}
	user, err := models.GetUser(u.db, &reqUser)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unable to login"})
		return
	}

	expiry := time.Now().Add(time.Hour * 24 * 7)
	expiryInterval := 60 * 60 * 24 * 7

	token, tokenErr := createToken(user.ID, expiry, u.secret)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to login"})
		return
	}

	c.SetCookie("token", token, expiryInterval, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"email":      user.Email,
		"resourceId": user.ID,
		"token":      token,
		"message":    "User logged in successfully!",
	})
}
