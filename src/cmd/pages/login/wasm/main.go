//go:generate tinygo build -o ../module.wasm -no-debug -panic=trap -scheduler=none main.go
package main

import (
	"ponglehub.co.uk/book-planner-go/src/pkg/web/wasm"
)

func main() {}

//export validate
func validate() {
	doc := wasm.GetDocument()

	value := doc.GetValue("username")
	password := doc.GetValue("password")

	if len(value) > 0 && len(password) > 0 {
		doc.SetDisabled("submit", false)
	} else {
		doc.SetDisabled("submit", true)
	}
}
