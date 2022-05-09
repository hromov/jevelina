package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
)

type UserTest struct {
	name                 string
	expectedResponseCode int
}

func TestUserHandler(t *testing.T) {
	tests := []UserTest{
		{
			name:                 "InitUsers",
			expectedResponseCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/users", nil)
			response := httptest.NewRecorder()

			_, err := cdb.OpenAndInitTest()
			if err != nil {
				log.Fatalf("Cant open and init data base error: %s", err.Error())
			}

			UsersHandler(response, request)

			if response.Result().StatusCode != test.expectedResponseCode {
				t.Errorf("Expected status code is: %d, but it's: %d", response.Result().StatusCode, test.expectedResponseCode)
			}

			var users []models.User
			// log.Println(response.Body.String())
			if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
				t.Errorf("Users JSON error: %s", err.Error())
			}
			if users == nil {
				t.Errorf("Users are nil")
			}
			initUsers := cdb.GetInitUsers()
			if len(users) < len(initUsers) {
				t.Errorf("It's expected to have at least %d user in DB, but it's only %d back", len(initUsers), len(users))
			}
		})
	}
}
