//go:generate tinygo build -o ../module.wasm -no-debug -panic=trap -scheduler=none main.go
package main

import (
	"strconv"
	"syscall/js"
)

func main() {
}

//export multiply
func multiply(x, y int) {
	doc := js.Global().Get("document")
	body := doc.Call("getElementById", "target")
	body.Set("innerHTML", strconv.FormatInt(int64(x*y), 10))
}
