package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/tests/validate"
)

func TestGetUserIntegration(t *testing.T) {
	type getUserSpec struct {
		user     string
		response *client.GetUserResponse
		err      string
	}

	test(t, []testSpec[getUserSpec]{
		{
			name: "Empty",
			spec: getUserSpec{
				user: "myuser",
				err:  "failed with status code 404",
			},
		},
		{
			name: "Found",
			existing: []database.User{
				{Name: "myuser", Password: "Password1?"},
			},
			spec: getUserSpec{
				user: "myuser",
				response: &client.GetUserResponse{
					Username: "myuser",
				},
			},
		},
		{
			name: "Wrong",
			existing: []database.User{
				{Name: "otheruser", Password: "Password1?"},
			},
			spec: getUserSpec{
				user: "myuser",
				err:  "failed with status code 404",
			},
		},
		{
			name: "Picked",
			existing: []database.User{
				{Name: "myuser", Password: "Password1?"},
				{Name: "youruser", Password: "Password2!"},
				{Name: "diffuser", Password: "Password3@"},
			},
			spec: getUserSpec{
				user: "youruser",
				response: &client.GetUserResponse{
					Username: "youruser",
				},
			},
		},
	}, func(u *testing.T, c *client.Client, spec getUserSpec) {
		response, err := c.GetUser(context.TODO(), spec.user)

		if spec.response != nil {
			validate.Equal(u, spec.response, response)
		} else {
			assert.Nil(u, response)
		}

		if spec.err != "" {
			assert.EqualError(u, err, spec.err)
		} else {
			assert.NoError(u, err)
		}
	})
}
