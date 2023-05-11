package footer

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/component"
)

//go:embed footer.html
var content string

func Get() component.Component {
	return component.Component{
		Name:     "components.footer",
		Template: content,
	}
}
