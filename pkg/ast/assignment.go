package ast

import (
	"bytes"
	"context"
	"reflect"
	"strings"

	"github.com/jippi/dottie/pkg/template"
	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
	slogctx "github.com/veqryn/slog-context"
)

type Assignment struct {
	Complete     bool                         `json:"complete"`     // The key/value had no value/content after the "=" sign
	Enabled      bool                         `json:"enabled"`      // The assignment was enabled out (#KEY=VALUE)
	Interpolated string                       `json:"interpolated"` // Value of the key (after interpolation)
	Literal      string                       `json:"literal"`      // Value of the key (right hand side of the "=" sign)
	Name         string                       `json:"key"`          // Name of the key (left hand side of the "=" sign)
	Quote        token.Quote                  `json:"quote"`        // The style of quotes used for the assignment
	Position     Position                     `json:"position"`     // Information about position of the assignment in the file
	Comments     []*Comment                   `json:"comments"`     // Comments attached to the assignment (e.g. doc block before it)
	Dependencies map[string]template.Variable `json:"dependencies"` // Assignments that this assignment depends on
	Dependents   map[string]*Assignment       `json:"dependents"`   // Assignments dependents on this assignment
	Group        *Group                       `json:"-"`            // The (optional) group this assignment belongs to
}

func (a *Assignment) statementNode() {}

func (a *Assignment) Initialize() {
	if dependencies := template.ExtractVariables(a.Literal); len(dependencies) > 0 {
		a.Dependencies = dependencies
	}
}

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

func (a *Assignment) SetLiteral(ctx context.Context, in string) {
	slogctx.Debug(ctx, "Assignment.SetLiteral() input", tui.StringDump("literal", in))

	a.Literal = token.Escape(ctx, a.Literal, a.Quote)
	a.Interpolated = a.Literal

	slogctx.Debug(ctx, "Assignment.SetLiteral() output", tui.StringDump("literal", a.Literal))
}

func (a *Assignment) Unquote(ctx context.Context) (string, error) {
	slogctx.Debug(ctx, "Assignment.Unquote() input", tui.StringDump("literal", a.Literal))

	str, err := token.Unescape(ctx, a.Literal, a.Quote)
	if err != nil {
		slogctx.Error(ctx, "failed to unquote string", tui.StringDump("literal", a.Literal))

		return "", err
	}

	slogctx.Debug(ctx, "Assignment.Unquote() output", tui.StringDump("literal", str))

	return str, nil
}
