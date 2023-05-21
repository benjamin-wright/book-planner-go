package integration

import (
	"flag"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/pkg/tests/cockroach"
	mockNats "ponglehub.co.uk/book-planner-go/src/pkg/tests/nats"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		m.Run()
		return
	}

	if testing.Verbose() {
		logger, _ := zap.NewDevelopment()
		zap.ReplaceGlobals(logger)
	} else {
		logger := zap.NewNop()
		zap.ReplaceGlobals(logger)
	}

	closeCockroach := cockroach.Run("cockroach", 26257)
	defer closeCockroach()

	closeNats := mockNats.Run("nats", 4222)
	defer closeNats()

	cockroach.Migrate("../migrations/books-1.sql")

	m.Run()
}

type testSpec[T any] struct {
	name     string
	existing []database.Book
	spec     T
}

func test[T any](t *testing.T, specs []testSpec[T], f func(t *testing.T, cli *database.Client, cn *nats.EncodedConn, spec T)) {
	if testing.Short() {
		t.SkipNow()
	}

	cli, err := database.New()
	if !assert.NoError(t, err) {
		return
	}

	nc := mockNats.Connect(t)
	defer nc.Close()

	for _, spec := range specs {
		t.Run(spec.name, func(u *testing.T) {
			if !assert.NoError(u, cli.DeleteAllBooks()) {
				return
			}

			for _, book := range spec.existing {
				if err := cli.AddBook(book); !assert.NoError(u, err) {
					return
				}
			}

			f(u, cli, nc, spec.spec)
		})
	}
}
