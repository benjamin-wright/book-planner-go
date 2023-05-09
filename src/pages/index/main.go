package main

import (
	_ "embed"
	"net/http"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/renderer"
)

//go:embed index.html
var content string

type Context struct {
	Title string
}

func main() {
	renderer := renderer.New(content)

	framework.Run(func(w http.ResponseWriter, r *http.Request) {
		zap.S().Infof("Got request")
		renderer.Execute(w, Context{Title: "my-page"})
	})
}
