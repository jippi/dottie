package ast

import "reflect"

type Newline struct {
	Blank      bool
	LineNumber int
	Group      *Group
}

func (n *Newline) Is(other Statement) bool {
	return reflect.TypeOf(n) == reflect.TypeOf(other)
}

func (n *Newline) BelongsToGroup(config RenderSettings) bool {
	return n.Group == nil || n.Group.BelongsToGroup(config)
}

func (n *Newline) ShouldRender(config RenderSettings) bool {
	return config.WithBlankLines()
}

func (n *Newline) Render(config RenderSettings) string {
	return "\n"
}

func (n *Newline) statementNode() {
}
