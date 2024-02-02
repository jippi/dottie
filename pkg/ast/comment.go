package ast

import (
	"bytes"
	"reflect"

	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
)

// Comment node represents a comment statement.
type Comment struct {
	Value      string            `json:"value"`      // The actual comment value
	Annotation *token.Annotation `json:"annotation"` // If the comment was detected to be an annotation
	Group      *Group            `json:"-"`          // The (optional) group the comment belongs to
	Position   Position          `json:"position"`   // Information about position of the assignment in the file
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

	var buf bytes.Buffer

	if config.WithColors() {
		out := tui.Theme.Success.Printer(tui.RendererWithTTY(&buf))
		if c.Annotation != nil {
			out.Print("# ")
			out.ApplyStyle(tui.Bold).Print("@", c.Annotation.Key)
			out.Print(" ")
			out.Println(c.Annotation.Value)
		} else {
			out.Println(c.Value)
		}

		return buf.String()
	}

	buf.WriteString(c.Value)
	buf.WriteString("\n")

	return buf.String()
}

func (c *Comment) statementNode() {
}

func (c *Comment) String() string {
	return c.Value
}
