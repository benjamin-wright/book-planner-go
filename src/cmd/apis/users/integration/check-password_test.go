package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/tests/validate"
)

func TestCheckPasswordIntegration(t *testing.T) {
	type checkPasswordSpec struct {
		username string
		password string
		response *client.CheckPasswordResponse
		ok       bool
		err      string
	}

	test(t, []testSpec[checkPasswordSpec]{
		{
			name: "Missing",
			spec: checkPasswordSpec{
				username: "myuser",
				password: "Password1!",
			},
		},
		{
			name: "Bad password",
			existing: []database.User{
				{Name: "my-user", Password: "Password3$"},
			},
			spec: checkPasswordSpec{
				username: "my-user",
				password: "Password1!",
			},
		},
		{
			name: "Success",
			existing: []database.User{
				{Name: "my-user", Password: "Password3$"},
			},
			spec: checkPasswordSpec{
				username: "my-user",
				password: "Password3$",
				response: &client.CheckPasswordResponse{
					ID:       "a uuid",
					Username: "my-user",
				},
				ok: true,
			},
		},
	}, func(u *testing.T, c *client.Client, spec checkPasswordSpec) {
		response, ok, err := c.CheckPassword(context.TODO(), spec.username, spec.password)

		if spec.err != "" {
			assert.EqualError(u, err, spec.err)
		} else {
			assert.NoError(u, err)
		}

		assert.Equal(u, spec.ok, ok)

		if spec.response != nil {
			validate.Equal(u, spec.response, response)
		} else {
			assert.Nil(u, response)
		}
	})
}
