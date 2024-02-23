package ast

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jippi/dottie/pkg/template"
	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
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

func (a *Assignment) SetLiteral(in string) {
	fmt.Printf("SetLiteral.input.string  >%s<\n", in)
	fmt.Printf("SetLiteral.input.unicode >%U<\n", []rune(in))

	a.Literal = tui.Quote(a.Literal)

	// val, err := tui.Unquote(in, '"', true)
	// if err != nil {
	// 	panic(err)
	// }

	// a.Literal = val

	fmt.Printf("EscapeString: out string >%s<\n", a.Literal)
	fmt.Printf("EscapeString: out unicode >%U<\n", []rune(a.Literal))

	a.Interpolated = a.Literal
}

func (a *Assignment) Unquote() string {
	fmt.Println("Unquote.input", fmt.Sprintf(">%q<", a.Literal))

	// str := tui.Quote(a.Literal)
	str, err := tui.Unquote(a.Literal, '"', true)
	if err != nil {
		panic(err)
	}

	newstr, err := strconv.Unquote("\"" + a.Literal + "\"")
	fmt.Println("Unquote.strconv.Unquote", fmt.Sprintf(">%q<", newstr), err)

	fmt.Println("Unquote.output", fmt.Sprintf(">%q<", str))

	return str
}
