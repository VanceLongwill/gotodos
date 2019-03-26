package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"strings"
	"testing"
)

type userTest struct {
	user     User
	expected string
}

func TestSerializeUser(t *testing.T) {
	t.Log(`should serialize user correctly`)
	userTests := make([]userTest, 2)

	user := User{
		ID:        1,
		FirstName: sql.NullString{"John", true},
		LastName:  sql.NullString{"Smith", true},
		Email:     "a@b.com",
		Password:  "password",
	}

	userTests[0].user = user
	userTests[0].expected = `{"email":"a@b.com","firstName":"John","id":1,"lastName":"Smith"}`
	user.FirstName.Valid = false
	user.LastName.Valid = false
	userTests[1].user = user
	userTests[1].expected = `{"email":"a@b.com","id":1}`

	for _, test := range userTests {
		serialUser := test.user.Serialize()

		data, err := json.Marshal(serialUser)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		received := string(data)

		if strings.Compare(test.expected, received) != 0 {
			fmt.Println(" === EXPECTED ===\n", test.expected)
			fmt.Println(" === RECEIVED ===\n", received)
			t.Fail()
		}
	}
}

var userTableRows = []string{"id", "first_name", "last_name", "email", "password_hash"}

func TestGetUser(t *testing.T) {
	t.Log(`Should get a single application user`)
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	userEmail := "a@b.com"

	rows := sqlmock.NewRows(userTableRows).
		AddRow(1, "John", "Smith", userEmail, "password")

	mock.ExpectQuery(`SELECT \* FROM users WHERE.+`).
		WithArgs(userEmail).
		WillReturnRows(rows)

	db := DB{mockDB}

	user, err := db.GetUser(userEmail)

	t.Log(user.Serialize())

	if err != nil {
		t.Errorf("failed while finding user")
		t.Fail()
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	t.Log(`should create a new application user`)

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer mockDB.Close()

	user := User{
		Email:     "a@b.com",
		FirstName: sql.NullString{"John", true},
		LastName:  sql.NullString{"Smith", true},
		Password:  "hashed_password",
	}

	mock.ExpectQuery(`INSERT INTO users.+`).
		WithArgs(user.Email, user.FirstName.String, user.LastName.String, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	db := DB{mockDB}

	registeredUser, err := db.CreateUser(&user)

	if err != nil {
		t.Errorf("failed while finding user %s", err.Error())
		t.Fail()
	}

	t.Log(registeredUser.Serialize())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
