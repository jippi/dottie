package ast

import (
	"reflect"
	"strings"

	"github.com/jippi/dottie/pkg/token"
)

// Comment node represents a comment statement.
type Comment struct {
	Value      string            `json:"value"`      // The actual comment value
	Annotation *token.Annotation `json:"annotation"` // If the comment was detected to be an annotation
	Group      *Group            `json:"-"`          // The (optional) group the comment belongs to
	Position   Position          `json:"position"`   // Information about position of the assignment in the file
}

func NewCommentsFromSlice(commentsSlice []string) []*Comment {
	if len(commentsSlice) == 0 {
		return nil
	}

	comments := make([]*Comment, len(commentsSlice))

	for i, comment := range commentsSlice {
		comments[i] = NewComment(comment)
	}

	return comments
}

func NewComment(value string) *Comment {
	return &Comment{
		Value: "# " + value,
	}
}

func (c *Comment) Is(other Statement) bool {
	if c == nil || other == nil {
		return false
	}

	return c.Type() == other.Type()
}

func (c *Comment) Type() string {
	if c == nil {
		return "<nil>Comment"
	}

	return reflect.TypeOf(c).String()
}

func (c *Comment) BelongsToGroup(name string) bool {
	if c.Group == nil && len(name) > 0 {
		return false
	}

	return c.Group == nil || c.Group.BelongsToGroup(name)
}

func (c *Comment) statementNode() {
}

func (c Comment) String() string {
	return c.Value
}

func (c Comment) CleanString() string {
	return strings.TrimPrefix(c.Value, "# ")
}
