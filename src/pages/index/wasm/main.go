//go:generate tinygo build -o module.wasm -no-debug -panic=trap -scheduler=none -gc=leaking main.go

package main

func main() {
	println("Hello World!")
}
