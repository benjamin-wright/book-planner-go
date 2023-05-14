package main

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed styles.css
var styles []byte

func main() {
	runtime.RunFileServer([]runtime.ServeFile{
		{Path: "styles.css", Data: styles, MimeType: "text/css"},
	})
}
