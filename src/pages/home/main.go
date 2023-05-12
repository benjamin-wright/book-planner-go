package main

import (
	_ "embed"
	"net/http"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

type Context struct {
}

func main() {
	runtime.Run(runtime.ServerOptions{
		Template: content,
		Title:    "Book Planner",
		Handler: func(r *http.Request) any {
			return Context{}
		},
	})
}
