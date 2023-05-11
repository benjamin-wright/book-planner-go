package wasm

import (
	_ "embed"
)

//go:embed module.wasm
var module []byte

func Module() []byte {
	return module
}
