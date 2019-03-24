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
)

func createToken(userID uint, expireTime time.Time, secret []byte) (string, error) {
	type Claims = middleware.UserClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	})
	return token.SignedString(secret)
}

// User is a type alias for convenience
type User = models.User

// UserStore represents the database layer for Users
type UserStore interface {
	DBGetUser
	DBCreateUser
}

// DBCreateUser represents the part of the datalayer responsible for creation of users
type DBCreateUser interface {
	CreateUser(u *User) (*User, error)
}

// RegisterUser returns a function which handles requests to create a new application user
func RegisterUser(db DBCreateUser, secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// RegisterRequest specifies the request body shape for registering a new application user
		type RegisterRequest struct {
			Email     string `json: "email" binding: "required"`
			Password  string `json: "password" binding: "required"`
			FirstName string `json: "firstName" binding: "required"`
			LastName  string `json: "lastName" binding: "required"`
		}
		var body RegisterRequest
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Bad request",
			})
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
			FirstName: sql.NullString{body.FirstName, true},
			LastName:  sql.NullString{body.LastName, true},
		}

		newUser, err := db.CreateUser(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Unable to register user"})
			return
		}
		expiryInterval := 60 * 60 * 24 * 7
		expiry := time.Now().Add(time.Hour * 24 * 7)

		token, tokenErr := createToken(user.ID, expiry, secret)
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
}

// DBGetUser represents the part of the datalayer responsible for getting a single user
type DBGetUser interface {
	GetUser(email string) (*User, error)
}

// LoginUser returns a function which handles requests to login application users
func LoginUser(db DBGetUser, secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// LoginRequest specifies the request body shape for logging in an application user
		type LoginRequest struct {
			Email    string `json: "email" binding: "required"`
			Password string `json: "password" binding: "required"`
		}
		var body LoginRequest
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Bad request",
			})
			return
		}

		user, err := db.GetUser(body.Email)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid password"})
			return
		}

		expiry := time.Now().Add(time.Hour * 24 * 7)
		expiryInterval := 60 * 60 * 24 * 7

		token, tokenErr := createToken(user.ID, expiry, secret)
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
}
