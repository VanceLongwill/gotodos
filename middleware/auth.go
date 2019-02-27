package middleware

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// UserClaims is used for creating and parsing jwts
type UserClaims struct {
	ID                 uint `json: "id"`
	jwt.StandardClaims      // includes ExpiresAt
}

// InvalidTokenError is used to distinguish between different errors in Authorization middleware
var InvalidTokenError error

func validateToken(inToken string, claims *UserClaims) error {
	token, err := jwt.ParseWithClaims(inToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(inToken), nil
	})

	if err != nil {
		return err
	}
	if !token.Valid {
		return InvalidTokenError
	}
	return nil
}

// Authorize blocks unauthorized requests
func Authorize(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, cookieErr := c.Cookie("token")
		if cookieErr != nil { // No token cookie provided
			authorization := c.Request.Header.Get("Authorization")
			if len(authorization) == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Authorized routes require cookie token or Authorization header",
				})
				return
			}
			authorizationParts := strings.Split(authorization, " ")
			if len(authorizationParts) != 2 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusBadRequest,
					"message": "Authorization header should be in the format `Bearer $TOKEN`",
				})
				return
			}
			tokenString = authorizationParts[1]
		}

		var claims UserClaims
		if err := validateToken(tokenString, &claims); err != nil {
			if err == jwt.ErrSignatureInvalid || err == InvalidTokenError {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Invalid token",
				})
			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusBadRequest,
					"message": "Unable to parse token",
				})
			}
			return
		}

		// JWT is valid, proceed
		c.Set("userID", claims.ID)
		c.Next()
	}
}
