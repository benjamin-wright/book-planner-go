package main

import (
	_ "embed"
	"net/http"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

type Context struct {
}

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	server := runtime.NewServer(content, "Book Planner")
	server.Run(func(r *http.Request) any {
		return Context{}
	})
}
