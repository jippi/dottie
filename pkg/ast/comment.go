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

func (c *Comment) Is(other Statement) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(other)
}

func (c *Comment) BelongsToGroup(config RenderSettings) bool {
	if c.Group == nil && len(config.FilterGroup) > 0 {
		return false
	}

	return c.Group == nil || c.Group.BelongsToGroup(config)
}

func (c *Comment) ShouldRender(config RenderSettings) bool {
	return config.WithComments() && c.BelongsToGroup(config)
}

func (c *Comment) Render(config RenderSettings) string {
	var buff bytes.Buffer

	buff.WriteString(c.String())
	buff.WriteString("\n")

	return buff.String()
}

func (c *Comment) statementNode() {
}

func (c *Comment) String() string {
	return c.Value
}
