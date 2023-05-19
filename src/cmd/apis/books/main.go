package main

import (
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/internal/server"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	api.Init()

	c, err := database.New()
	if err != nil {
		zap.S().Fatalf("Failed to create postgres client: %+v", err)
	}

	api.Run(server.Router(c))
}
