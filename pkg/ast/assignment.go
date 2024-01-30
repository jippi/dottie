package ast

import (
	"bytes"
	"fmt"
	"reflect"

	"dotfedi/pkg/token"
)

type Assignment struct {
	Key               string          `json:"key"`
	Value             string          `json:"value"`
	CompleteStatement bool            `json:"complete"`
	Active            bool            `json:"commented"`
	QuoteType         token.QuoteType `json:"quote"`
	Group             *Group          `json:"-"`
	Position          Position        `json:"position"`
	Comments          []*Comment      `json:"comments"`
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
		a.QuoteType = token.DoubleQuotes
	case "'", "single":
		a.QuoteType = token.SingleQuotes
	case "none":
		a.QuoteType = token.NoQuotes
	}
}

func (a *Assignment) Assignment() string {
	if a.QuoteType == token.NoQuotes {
		return fmt.Sprintf("%s=%s", a.Key, a.Value)
	}

	return fmt.Sprintf("%s=%s%s%s", a.Key, a.QuoteType, a.Value, a.QuoteType)
}
