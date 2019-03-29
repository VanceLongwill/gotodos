package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var secret = []byte("some secret")

type mockCreateUser struct{}

func (db mockCreateUser) CreateUser(u *User) (*User, error) {
	return u, nil
}

type mockGetUser struct {
	password string
}

func (db mockGetUser) GetUser(email string) (*User, error) {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(db.password), 12)
	return &User{Email: email, Password: string(hashBytes)}, nil
}

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []mock{
		{
			`{}`,
			http.StatusBadRequest,
		},
		{
			`{"some": "field"}`,
			http.StatusBadRequest,
		},
		{
			`{"email": "a@b.com", "password": "password", "firstName": "Bob", "lastName": "Smith"}`,
			http.StatusCreated,
		},
	}

	for _, test := range tests {
		recorder := httptest.NewRecorder() // implements http.ResponseWriter
		mockContext, _ := gin.CreateTestContext(recorder)
		userID := uint(11)
		mockContext.Set("userID", userID)
		req, _ := http.NewRequest("POST", "http://example.com/",
			bytes.NewBuffer([]byte(test.json)))
		mockContext.Request = req
		db := mockCreateUser{}
		RegisterUser(db, secret)(mockContext)

		if recorder.Code != test.expectedCode {
			t.Errorf("Expected status code %d but received %d",
				test.expectedCode, recorder.Code)
			t.Log(recorder.Body)
		}
	}
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []mock{
		{
			`{}`,
			http.StatusBadRequest,
		},
		{
			`{"email": "field"}`,
			http.StatusBadRequest,
		},
		{
			`{"email": "a@b.com", "password": "password"}`,
			http.StatusOK,
		},
		{
			`{"email": "a@b.com", "password": "incorrect_password"}`,
			http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		recorder := httptest.NewRecorder() // implements http.ResponseWriter
		mockContext, _ := gin.CreateTestContext(recorder)
		userID := uint(11)
		mockContext.Set("userID", userID)
		req, _ := http.NewRequest("POST", "http://example.com/",
			bytes.NewBuffer([]byte(test.json)))
		mockContext.Request = req

		db := mockGetUser{"password"}
		LoginUser(db, secret)(mockContext)

		if recorder.Code != test.expectedCode {
			t.Errorf("Expected status code %d but received %d",
				test.expectedCode, recorder.Code)
			t.Log(recorder.Body)
		}
	}
}
