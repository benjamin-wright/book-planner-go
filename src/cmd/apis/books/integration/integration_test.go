package integration

import (
	"flag"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/internal/server"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/tests/cockroach"
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
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		logger := zap.NewNop()
		zap.ReplaceGlobals(logger)
	}

	close := cockroach.Run("cockroach", 26257)
	defer close()

	cockroach.Migrate("../migrations/books-1.sql")

	m.Run()
}

type testSpec[T any] struct {
	name     string
	existing []database.Book
	spec     T
}

func test[T any](t *testing.T, specs []testSpec[T], f func(t *testing.T, c *client.Client, spec T)) {
	if testing.Short() {
		t.SkipNow()
	}

	cli, err := database.New()
	if !assert.NoError(t, err) {
		return
	}

	r := server.Router(cli)
	srv := httptest.NewServer(r.Handler())
	c := client.New(srv.URL)

	for _, spec := range specs {
		t.Run(spec.name, func(u *testing.T) {
			if !assert.NoError(u, cli.DeleteAllBooks()) {
				return
			}

			for _, book := range spec.existing {
				if err := cli.AddBook(book); !assert.NoError(t, err) {
					return
				}
			}

			f(u, c, spec.spec)
		})
	}
}
