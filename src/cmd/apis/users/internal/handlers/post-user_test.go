package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/types"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
)

func TestPostUserIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	cli, err := database.New()
	if !assert.NoError(t, err) {
		return
	}

	handler := PostUser(cli)

	for _, test := range []struct {
		name     string
		request  client.AddUserRequest
		existing []types.User
		status   int
	}{
		{
			name:    "Short password",
			request: client.AddUserRequest{Username: "myuser", Password: "hi"},
			status:  http.StatusBadRequest,
		},
		{
			name:    "Simple password",
			request: client.AddUserRequest{Username: "myuser", Password: "longbutsimple"},
			status:  http.StatusBadRequest,
		},
		{
			name:    "Success",
			request: client.AddUserRequest{Username: "myuser", Password: "Password1?"},
			status:  http.StatusCreated,
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			if !assert.NoError(u, cli.DeleteAllUsers()) {
				return
			}

			for _, user := range test.existing {
				if err := cli.AddUser(user.Name, user.Password); !assert.NoError(t, err) {
					return
				}
			}

			r := handler.TestHandler()
			w := httptest.NewRecorder()

			body, err := json.Marshal(&test.request)
			if !assert.NoError(t, err) {
				return
			}

			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
			r.ServeHTTP(w, req)

			assert.Equal(u, test.status, w.Code)
		})
	}
}
