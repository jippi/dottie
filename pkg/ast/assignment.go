package ast

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg/token"
)

type Assignment struct {
	Name         string      `json:"key"`       // Name of the key (left hand side of the "=" sign)
	Literal      string      `json:"literal"`   // Value of the key (right hand side of the "=" sign)
	Interpolated string      `json:"value"`     // Value of the key (after interpolation)
	Complete     bool        `json:"complete"`  // The key/value had no value/content after the "=" sign
	Active       bool        `json:"commented"` // The assignment was commented out (#KEY=VALUE)
	Quote        token.Quote `json:"quote"`     // The style of quotes used for the assignment
	Group        *Group      `json:"-"`         // The (optional) group this assignment belongs to
	Comments     []*Comment  `json:"comments"`  // Comments attached to the assignment (e.g. doc block before it)
	Position     Position    `json:"position"`  // Information about position of the assignment in the file
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

func (a *Assignment) Documentation(withoutPrefix bool) string {
	var buff bytes.Buffer

	for _, c := range a.Comments {
		val := c.Value

		if withoutPrefix {
			val = strings.TrimPrefix(val, "# ")
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

func (a *Assignment) IsValid() error {
	if !a.Quote.Valid() {
		return fmt.Errorf("invalid quote-style")
	}

	return validator.
		New(validator.WithRequiredStructEnabled()).
		Var(a.Interpolated, a.ValidationRules())
}

func (a *Assignment) Disable() {
	a.Active = false
}

func (a *Assignment) Enable() {
	a.Active = true
}
