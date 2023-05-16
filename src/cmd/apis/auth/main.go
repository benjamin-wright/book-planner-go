package main

import (
	"os"

	"ponglehub.co.uk/book-planner-go/src/cmd/apis/auth/internal/handlers"
	tokensApi "ponglehub.co.uk/book-planner-go/src/cmd/apis/tokens/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	loginURL := os.Getenv("LOGIN_URL")
	tokens := tokensApi.New(os.Getenv("TOKENS_API_URL"))

	api.Run(api.RunOptions{
		Handlers: []api.Handler{
			handlers.Verify(tokens, loginURL),
		},
	})
}
