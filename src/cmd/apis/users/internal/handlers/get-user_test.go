package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/tests/validate"
)

func TestGetUserIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	cli, err := database.New()
	if !assert.NoError(t, err) {
		return
	}

	handler := GetUser(cli)

	type testUser struct {
		name     string
		password string
	}

	for _, test := range []struct {
		name     string
		user     string
		existing []testUser
		status   int
		response *client.GetUserResponse
	}{
		{
			name:   "Empty",
			user:   "myuser",
			status: http.StatusNotFound,
		},
		{
			name: "Found",
			user: "myuser",
			existing: []testUser{
				{name: "myuser", password: "Password1?"},
			},
			status: http.StatusOK,
			response: &client.GetUserResponse{
				Username: "myuser",
			},
		},
		{
			name: "Wrong",
			user: "myuser",
			existing: []testUser{
				{name: "otheruser", password: "Password1?"},
			},
			status: http.StatusNotFound,
		},
		{
			name: "Picked",
			user: "youruser",
			existing: []testUser{
				{name: "myuser", password: "Password1?"},
				{name: "youruser", password: "Password2!"},
				{name: "diffuser", password: "Password3@"},
			},
			status: http.StatusOK,
			response: &client.GetUserResponse{
				Username: "youruser",
			},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			if !assert.NoError(u, cli.DeleteAllUsers()) {
				return
			}

			for _, user := range test.existing {
				if err := cli.AddUser(user.name, user.password); !assert.NoError(t, err) {
					return
				}
			}

			r := handler.TestHandler()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/"+test.user, nil)
			r.ServeHTTP(w, req)

			assert.Equal(u, test.status, w.Code)

			if test.response != nil {
				var result client.GetUserResponse
				if !assert.NoError(u, json.Unmarshal(w.Body.Bytes(), &result)) {
					return
				}

				validate.Equal(u, *test.response, result)
			}
		})
	}
}
