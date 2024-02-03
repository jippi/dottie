package ast

import (
	"reflect"
)

type Newline struct {
	Blank    bool   `json:"blank"`
	Group    *Group `json:"-"`
	Position Position
}

func (n *Newline) Is(other Statement) bool {
	return reflect.TypeOf(n) == reflect.TypeOf(other)
}

func (n *Newline) BelongsToGroup(config RenderSettings) bool {
	return n.Group == nil || n.Group.BelongsToGroup(config)
}

func (n *Newline) Render(config RenderSettings) string {
	if !config.WithBlankLines() {
		return ""
	}

	return "\n"
}

func (n *Newline) statementNode() {
}
