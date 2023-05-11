package header

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/component"
)

//go:embed header.html
var content string

func Get() component.Component {
	return component.Component{
		Name:     "components.header",
		Template: content,
	}
}
