package main

import (
	_ "embed"
	"net/http"
	"os"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

type Context struct {
	RegisterURL string
}

func main() {
	registerURL := os.Getenv("REGISTER_URL")

	runtime.Run(runtime.ServerOptions{
		Template: content,
		Title:    "Book Planner",
		PageHandler: func(r *http.Request) any {
			return Context{
				RegisterURL: registerURL,
			}
		},
	})
}
