package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/handlers"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/types"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
)

func TestPutPasswordIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	cli, err := database.New()
	if !assert.NoError(t, err) {
		return
	}

	handler := handlers.PutPassword(cli)

	for _, test := range []struct {
		name     string
		url      string
		request  client.CheckPasswordRequest
		existing []types.User
		status   int
	}{
		{
			name:    "Missing",
			url:     "/my-user/password",
			request: client.CheckPasswordRequest{Password: "Password3$"},
			status:  http.StatusUnauthorized,
		},
		{
			name:    "Bad password",
			url:     "/my-user/password",
			request: client.CheckPasswordRequest{Password: "Password1!"},
			existing: []types.User{
				{Name: "my-user", Password: "Password3$"},
			},
			status: http.StatusUnauthorized,
		},
		{
			name:    "Success",
			url:     "/my-user/password",
			request: client.CheckPasswordRequest{Password: "Password3$"},
			existing: []types.User{
				{Name: "my-user", Password: "Password3$"},
			},
			status: http.StatusOK,
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

			r := handler.TestHandler(testing.Verbose())
			w := httptest.NewRecorder()

			body, err := json.Marshal(test.request)
			if !assert.NoError(t, err) {
				return
			}

			req, _ := http.NewRequest("PUT", test.url, bytes.NewBuffer(body))
			r.ServeHTTP(w, req)

			assert.Equal(u, test.status, w.Code)
		})
	}
}
