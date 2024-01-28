package ast

import (
	"bytes"
	"reflect"
)

// Comment node represents a comment statement.
type Comment struct {
	Value           string
	LineNumber      int
	Annotation      bool
	AnnotationKey   string
	AnnotationValue string
	Group           *Group
}

func NewComment(value string) *Comment {
	return &Comment{
		Value: "# " + value,
	}
}

func (s *Comment) Is(other Statement) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s *Comment) BelongsToGroup(config RenderSettings) bool {
	if s.Group == nil && len(config.FilterGroup) > 0 {
		return false
	}

	return s.Group == nil || s.Group.BelongsToGroup(config)
}

func (s *Comment) ShouldRender(config RenderSettings) bool {
	return config.WithComments() && s.BelongsToGroup(config)
}

func (s *Comment) Render(config RenderSettings) string {
	var buff bytes.Buffer

	buff.WriteString(s.String())
	buff.WriteString("\n")

	return buff.String()
}

func (s *Comment) statementNode() {
}

func (s *Comment) String() string {
	return s.Value
}
