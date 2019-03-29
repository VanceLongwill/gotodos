package middleware

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := []byte("secret")

	type headerTest struct {
		header string
		code   int
	}

	invalidHeaders := []headerTest{{"Bearer token", http.StatusBadRequest},
		{"Bearer ", http.StatusBadRequest},
		{"token", http.StatusBadRequest},
		{"", http.StatusUnauthorized},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		ID: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix(),
		},
	})

	tokenStr, _ := token.SignedString(secret)
	authString := fmt.Sprintf("Bearer %s", tokenStr)

	invalidHeaders = append(invalidHeaders,
		[]headerTest{{authString, http.StatusOK}}...)

	for _, header := range invalidHeaders {
		recorder := httptest.NewRecorder() // implements http.ResponseWriter
		mockContext, _ := gin.CreateTestContext(recorder)
		req, _ := http.NewRequest("GET", "/", nil)
		mockContext.Request = req
		mockContext.Request.Header.Add("Authorization", header.header)

		Authorize(secret)(mockContext)

		if recorder.Code != header.code {
			t.Errorf("Expected status code %d but received %d", header.code, recorder.Code)
			t.Log(recorder.Body)
			t.Fail()
			return
		}
	}
}
