package main

import (
	"os"

	"ponglehub.co.uk/book-planner-go/src/cmd/apis/tokens/internal/handlers"
	"ponglehub.co.uk/book-planner-go/src/internal/tokens"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	keyfile := tokens.Keyfile(os.Getenv("TOKEN_KEYFILE"))

	api.Run(api.Router(api.RunOptions{
		Handlers: []api.Handler{
			handlers.GetLoginToken(keyfile),
			handlers.ValidateLoginToken(keyfile),
		},
	}))
}
