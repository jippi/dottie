package ast

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/jippi/dottie/pkg/token"
)

type Assignment struct {
	Comments     []*Comment  `json:"comments"` // Comments attached to the assignment (e.g. doc block before it)
	Complete     bool        `json:"complete"` // The key/value had no value/content after the "=" sign
	Enabled      bool        `json:"enabled"`  // The assignment was enabled out (#KEY=VALUE)
	Group        *Group      `json:"-"`        // The (optional) group this assignment belongs to
	Interpolated string      `json:"value"`    // Value of the key (after interpolation)
	Literal      string      `json:"literal"`  // Value of the key (right hand side of the "=" sign)
	Name         string      `json:"key"`      // Name of the key (left hand side of the "=" sign)
	Position     Position    `json:"position"` // Information about position of the assignment in the file
	Quote        token.Quote `json:"quote"`    // The style of quotes used for the assignment
}

func (a *Assignment) statementNode() {}

func (a *Assignment) Is(other Statement) bool {
	if a == nil || other == nil {
		return false
	}

	return a.Type() == other.Type()
}

func (a *Assignment) Type() string {
	if a == nil {
		return "<nil>Assignment"
	}

	return reflect.TypeOf(a).String()
}

func (a *Assignment) BelongsToGroup(name string) bool {
	if a.Group == nil && len(name) > 0 {
		return false
	}

	return a.Group == nil || a.Group.BelongsToGroup(name)
}

func (a *Assignment) HasComments() bool {
	return len(a.Comments) > 0
}

func (a *Assignment) DocumentationSummary() string {
	if !a.HasComments() {
		return ""
	}

	return strings.TrimPrefix(a.Comments[0].String(), "#")
}

func (a *Assignment) Documentation(withoutPrefix bool) string {
	var buff bytes.Buffer

	for _, comment := range a.Comments {
		// Exclude annotations from documentation
		if comment.Annotation != nil {
			continue
		}

		val := comment.Value

		if withoutPrefix {
			val = strings.TrimPrefix(val, "#")
		}

		buff.WriteString(val)
		buff.WriteString("\n")
	}

	return buff.String()
}

func (a *Assignment) ValidationRules() string {
	for _, comment := range a.Comments {
		if comment.Annotation == nil {
			continue
		}

		if comment.Annotation.Key == "dottie/validate" {
			return comment.Annotation.Value
		}
	}

	return ""
}

func (a *Assignment) IsHidden() bool {
	for _, comment := range a.Comments {
		if comment.Annotation == nil {
			continue
		}

		return comment.Annotation.Key == "dottie/hidden"
	}

	return false
}

func (a *Assignment) Disable() {
	a.Enabled = false
}

func (a *Assignment) Enable() {
	a.Enabled = true
}

func (a *Assignment) CommentsSlice() []string {
	res := []string{}

	for _, comment := range a.Comments {
		res = append(res, comment.CleanString())
	}

	return res
}
