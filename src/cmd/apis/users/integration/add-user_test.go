package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
)

func TestAddUserIntegration(t *testing.T) {
	type addUserSpec struct {
		username string
		password string
		err      string
	}

	test(t, []testSpec[addUserSpec]{
		{
			name: "Short password",
			spec: addUserSpec{
				username: "myuser",
				password: "hi",
				err:      "failed with status code 400",
			},
		},
		{
			name: "Simple password",
			spec: addUserSpec{
				username: "myuser",
				password: "longbutsimple",
				err:      "failed with status code 400",
			},
		},
		{
			name: "Success",
			spec: addUserSpec{
				username: "myuser",
				password: "Password1?",
			},
		},
		{
			name: "Exists",
			existing: []database.User{
				{Name: "myuser", Password: "Password2!"},
			},
			spec: addUserSpec{
				username: "myuser",
				password: "Password1?",
				err:      client.ErrUserExists.Error(),
			},
		},
	}, func(u *testing.T, c *client.Client, spec addUserSpec) {
		err := c.AddUser(context.TODO(), spec.username, spec.password)

		if spec.err != "" {
			assert.EqualError(u, err, spec.err)
		} else {
			assert.Nil(u, err)
		}
	})
}
