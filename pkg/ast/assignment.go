package ast

import (
	"bytes"
	"fmt"
	"reflect"

	"dotfedi/pkg/token"
)

type Assignment struct {
	Name     string      `json:"key"`       // Name of the key (left hand side of the "=" sign)
	Value    string      `json:"value"`     // Value of the key (right hand side of the "=" sign)
	Complete bool        `json:"complete"`  // The key/value had no value/content after the "=" sign
	Active   bool        `json:"commented"` // The assignment was commented out (#KEY=VALUE)
	Quote    token.Quote `json:"quote"`     // The style of quotes used for the assignment
	Group    *Group      `json:"-"`         // The (optional) group this assignment belongs to
	Comments []*Comment  `json:"comments"`  // Comments attached to the assignment (e.g. doc block before it)
	Position Position    `json:"position"`  // Information about position of the assignment in the file
}

func (a *Assignment) statementNode() {}

func (a *Assignment) Is(other Statement) bool {
	if other == nil || a == nil {
		return false
	}

	return reflect.TypeOf(a) == reflect.TypeOf(other)
}

func (a *Assignment) BelongsToGroup(config RenderSettings) bool {
	return a.Group == nil || a.Group.BelongsToGroup(config)
}

func (a *Assignment) HasComments() bool {
	return len(a.Comments) > 0
}

func (a *Assignment) Render(config RenderSettings) string {
	if !config.Match(a) || !a.BelongsToGroup(config) {
		return ""
	}

	var buff bytes.Buffer

	if config.WithComments() {
		for _, c := range a.Comments {
			buff.WriteString(c.Value)
			buff.WriteString("\n")
		}
	}

	if a.Active {
		buff.WriteString("#")
	}

	buff.WriteString(a.Assignment())
	buff.WriteString("\n")

	return buff.String()
}

func (a *Assignment) SetQuote(in string) {
	switch in {
	case "\"", "double":
		a.Quote = token.DoubleQuotes
	case "'", "single":
		a.Quote = token.SingleQuotes
	case "none":
		a.Quote = token.NoQuotes
	}
}

func (a *Assignment) Assignment() string {
	if a.Quote == token.NoQuotes {
		return fmt.Sprintf("%s=%s", a.Name, a.Value)
	}

	return fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, a.Value, a.Quote)
}
