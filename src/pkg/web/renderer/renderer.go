package renderer

import (
	_ "embed"
	"html/template"
	"io"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/components/footer"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/components/header"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework"
)

//go:embed base.html
var baseContent string

type Renderer struct {
	t *template.Template
}

func New(content string, children ...framework.Component) Renderer {
	t := template.New("Page")

	pageComponent := framework.Component{
		Name:     "page",
		Template: content,
		Children: append([]framework.Component{
			header.Get(),
			footer.Get(),
		}, children...),
	}

	pageComponent.Parse(t)
	template.Must(t.Parse(baseContent))

	return Renderer{t: t}
}

func (r *Renderer) Execute(wr io.Writer, data any) {
	r.t.Execute(wr, data)
}
