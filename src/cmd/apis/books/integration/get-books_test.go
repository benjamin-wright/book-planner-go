package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/tests/validate"
)

func TestGetBooksIntegration(t *testing.T) {
	type getBooksSpec struct {
		response *client.GetBooksResponse
		err      string
	}

	user := uuid.New().String()
	other_user := uuid.New().String()

	test(
		t,
		[]testSpec[getBooksSpec]{
			{
				name: "empty",
				spec: getBooksSpec{
					response: &client.GetBooksResponse{
						Books: []client.Book{},
					},
				},
			},
			{
				name: "one book",
				existing: []database.Book{
					{
						UserID:  user,
						Name:    "my book",
						Summary: "my summary",
					},
					{
						UserID:  other_user,
						Name:    "someone else's book",
						Summary: "another summary",
					},
				},
				spec: getBooksSpec{
					response: &client.GetBooksResponse{
						Books: []client.Book{
							{
								ID:          "some id",
								Name:        "my book",
								Summary:     "my summary",
								CreatedTime: time.Now(),
							},
						},
					},
				},
			},
		},
		func(u *testing.T, c *client.Client, spec getBooksSpec) {
			response, err := c.GetBooks(user)

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
		},
	)
}
