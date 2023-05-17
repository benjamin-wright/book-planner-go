package main

import (
	_ "embed"
	"net/http"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

func main() {
	runtime.Run(runtime.ServerOptions{
		Template: content,
		Title:    "Book Planner",
		PageHandler: func(r *http.Request) any {
			zap.S().Infof("Serving home page for user: %s", r.Header.Get("X-Auth-User"))
			return nil
		},
	})
}
