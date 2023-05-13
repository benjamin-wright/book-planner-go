package main

import (
	_ "embed"

	"go.uber.org/zap"
)

//go:embed index.html
var content string

type Context struct {
	Title   string
	Widgets []string
}

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// server := runtime.NewServer(content, "my-page", widget.Get())
	// server.AddWASMModule("wasm", "/main.wasm", wasm.Module())
	// server.Run(func(r *http.Request) any {
	// 	zap.S().Infof("Got request")
	// 	return Context{
	// 		Title:   "my-page",
	// 		Widgets: []string{"thing1", "thing2", "thing3"},
	// 	}
	// })
}
