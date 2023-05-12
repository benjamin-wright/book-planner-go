package widget

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/component"
)

//go:embed widget.html
var content string

func Get() component.Component {
	return component.Component{
		Name:     "components.widget",
		Template: content,
	}
}
