package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/books/internal/handlers"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/pkg/types"
	mockNats "ponglehub.co.uk/book-planner-go/src/pkg/tests/nats"
)

func TestCreateBookIntegration(t *testing.T) {
	type createBooksSpec struct {
		event types.CreateBookEvent
		err   string
	}

	user := uuid.New().String()
	// other_user := uuid.New().String()

	test(
		t,
		[]testSpec[createBooksSpec]{
			{
				name: "Success",
				spec: createBooksSpec{
					event: types.CreateBookEvent{
						UserID:  user,
						Name:    "my-book",
						Summary: "a lovely book",
					},
				},
			},
		},
		func(u *testing.T, cli *database.Client, cn *nats.EncodedConn, spec createBooksSpec) {
			mockNats.Handle(t, cn, handlers.CreateBook(cli, "books.create").GetHandler())

			response := types.CreateBookResponseEvent{}
			err := cn.Request("books.create", spec.event, &response, time.Second*5)

			if spec.err != "" {
				assert.EqualError(u, err, spec.err)
			} else {
				assert.NoError(u, err)
			}
		},
	)
}
