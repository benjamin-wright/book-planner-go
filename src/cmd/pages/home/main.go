package main

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

func main() {
	runtime.Run(runtime.ServerOptions{
		Template: content,
		Title:    "Book Planner",
	})
}
