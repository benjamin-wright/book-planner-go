package main

import (
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/handlers"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	api.Init()

	cli, err := database.New()
	if err != nil {
		zap.S().Fatalf("Error getting database client: %+v", err)
	}

	api.Run(api.RunOptions{
		Handlers: []api.Handler{
			handlers.PostUser(cli),
			handlers.GetUser(cli),
			handlers.PutPassword(cli),
		},
	})
}
