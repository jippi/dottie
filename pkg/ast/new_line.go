package ast

import "reflect"

type Newline struct {
	Blank      bool
	LineNumber int
	Group      *Group
}

func (s *Newline) Is(other Statement) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s *Newline) BelongsToGroup(config RenderSettings) bool {
	return s.Group == nil || s.Group.BelongsToGroup(config)
}

func (s *Newline) ShouldRender(config RenderSettings) bool {
	return config.WithBlankLines()
}

func (s *Newline) Render(config RenderSettings) string {
	return "\n"
}

func (s *Newline) statementNode() {
}
