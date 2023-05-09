package header

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework"
)

//go:embed header.html
var content string

func Get() framework.Component {
	return framework.Component{
		Name:     "components.header",
		Template: content,
	}
}
