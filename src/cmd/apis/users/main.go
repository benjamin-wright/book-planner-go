package main

import (
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/server"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	api.Init()

	cli, err := database.New()
	if err != nil {
		zap.S().Fatalf("Error getting database client: %+v", err)
	}

	api.Run(server.Router(cli))
}
