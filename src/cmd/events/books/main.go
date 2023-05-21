package main

import (
	"os"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/books/internal/handlers"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/events"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	client, err := database.New()
	if err != nil {
		zap.S().Fatalf("Failed to create postgres client: %+v", err)
	}

	err = events.Run(handlers.Get(client, handlers.Config{
		CreateBooksSubject: os.Getenv("CREATE_BOOKS_SUBJECT"),
	}))
	if err != nil {
		zap.S().Fatalf("Running event loop failed: %+v", err)
	}
}
