package handlers

import (
	"flag"
	"testing"

	"ponglehub.co.uk/book-planner-go/src/pkg/tests/cockroach"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		m.Run()
		return
	}

	close := cockroach.Run("cockroach", 26257)
	defer close()

	cockroach.Migrate("../../migrations/users-1.sql")

	m.Run()
}
