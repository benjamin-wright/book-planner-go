package wasm

import (
	"errors"
	"syscall/js"
)

// body := doc.Call("getElementById", "target")
// body.Set("innerHTML", strconv.FormatInt(int64(x*y), 10))

type Document struct {
	doc js.Value
}

func GetDocument() *Document {
	doc := js.Global().Get("document")
	return &Document{
		doc: doc,
	}
}

func (d *Document) AttachInputHandler(id string, handler func(newValue string)) error {
	element := d.doc.Call("getElementById", id)
	if element.IsUndefined() {
		return errors.New("failed to find element" + id)
	}

	element.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) any {
		// for idx, value := range args {
		// 	println("Argument " + strconv.FormatInt(int64(idx), 10) + " is a " + value.Type().String())
		// }

		return nil
	}))

	return nil
}

func (d *Document) GetValue(id string) string {
	element := d.doc.Call("getElementById", id)
	if element.IsUndefined() {
		panic("failed to find element" + id)
	}

	value := element.Get("value").String()

	return value
}

func (d *Document) SetDisabled(id string, value bool) {
	element := d.doc.Call("getElementById", id)
	if element.IsUndefined() {
		panic("failed to find element" + id)
	}

	element.Set("disabled", value)
}
