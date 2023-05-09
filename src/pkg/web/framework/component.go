package framework

import (
	"html/template"
)

type Component struct {
	Name     string
	Template string
	Children []Component
}

func wrap(name string, content string) string {
	return "{{ define \"" + name + "\" }}" + content + "{{ end }}"
}

func (c *Component) AddChild(component Component) {
	c.Children = append(c.Children, component)
}

func (c *Component) Parse(t *template.Template) {
	for _, c := range c.Children {
		c.Parse(t)
	}

	template.Must(t.Parse(wrap(c.Name, c.Template)))
}
