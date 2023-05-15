//go:generate tinygo build -o ../module.wasm -no-debug -panic=trap -scheduler=none main.go
package main

import (
	"ponglehub.co.uk/book-planner-go/src/pkg/web/wasm/dom"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/wasm/validation"
)

func main() {}

//export validate
func validate() {
	doc := dom.GetDocument()

	password := doc.GetValue("password")
	confirm := doc.GetValue("confirm-password")

	if !validation.CheckPasswordComplexity(password) {
		doc.SetCustomValidity("password", "Your password must be at least 8 and no more than 50 letters long, and contain one uppercase, one lowercase, one numeric and one special letter")
	} else {
		doc.SetCustomValidity("password", "")
	}

	if password != confirm {
		doc.SetCustomValidity("confirm-password", "Your password and password confirmation should match!")
	} else {
		doc.SetCustomValidity("confirm-password", "")
	}
}
