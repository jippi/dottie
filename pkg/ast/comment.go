package ast

import (
	"bytes"
	"reflect"
)

// Comment node represents a comment statement.
type Comment struct {
	Value           string `json:"value"`
	LineNumber      int    `json:"line_number"`
	Annotation      bool   `json:"annotation"`
	AnnotationKey   string `json:"annotation_key,omitempty"`
	AnnotationValue string `json:"annotation_value,omitempty"`
	Group           *Group `json:"-"`
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

func (c *Comment) Render(config RenderSettings) string {
	if !config.WithComments() || !c.BelongsToGroup(config) {
		return ""
	}

	var buff bytes.Buffer

	buff.WriteString(c.Value)
	buff.WriteString("\n")

	return buff.String()
}

func (c *Comment) statementNode() {
}

func (c *Comment) String() string {
	return c.Value
}
