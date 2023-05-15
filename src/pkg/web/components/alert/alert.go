package alert

import (
	_ "embed"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/component"
)

//go:embed alert.html
var content string

func Get() component.Component {
	return component.Component{
		Name:     "components.alert",
		Template: content,
	}
}

func Lookup(err string, lookup map[string]string) string {
	if err == "" {
		return ""
	}

	if response, ok := lookup[err]; ok {
		return response
	} else {
		return "Something unexpected happened, please try again later"
	}
}
