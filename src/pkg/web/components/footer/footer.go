package footer

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework"
)

//go:embed footer.html
var content string

func Get() framework.Component {
	return framework.Component{
		Name:     "components.footer",
		Template: content,
	}
}
