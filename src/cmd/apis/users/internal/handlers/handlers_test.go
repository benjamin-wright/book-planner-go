package handlers_test

import (
	"flag"
	"testing"

	"go.uber.org/zap"
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
		logger := zap.NewNop()
		zap.ReplaceGlobals(logger)
	}

	close := cockroach.Run("cockroach", 26257)
	defer close()

	cockroach.Migrate("../../migrations/users-1.sql")

	m.Run()
}
