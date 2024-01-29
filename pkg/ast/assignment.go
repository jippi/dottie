package ast

import (
	"bytes"
	"fmt"
	"reflect"
)

const (
	SingleQuotes = '\''
	DoubleQuotes = '"'
	NoQuotes     = 0
)

type Assignment struct {
	Key       string     `json:"key"`
	Value     string     `json:"value"`
	Comments  []*Comment `json:"comments"`
	Group     *Group     `json:"-"`
	Commented bool       `json:"commented"`

	FirstLine  int  `json:"first_line"`
	LastLine   int  `json:"last_line"`
	LineNumber int  `json:"line_number"`
	Naked      bool `json:"naked"`
	Complete   bool `json:"complete"`
	Quoted     rune `json:"quote"`
}

func (a *Assignment) statementNode() {}

func (a *Assignment) Is(other Statement) bool {
	if other == nil {
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

	if a.Commented {
		buff.WriteString("#")
	}

	buff.WriteString(a.Assignment())
	buff.WriteString("\n")

	return buff.String()
}

func (a *Assignment) SetQuote(in string) {
	switch in {
	case "\"", "double":
		a.Quoted = DoubleQuotes
	case "'", "single":
		a.Quoted = SingleQuotes
	case "none":
		a.Quoted = NoQuotes
	}
}

func (a *Assignment) Assignment() string {
	if a.Quoted == NoQuotes {
		return fmt.Sprintf("%s=%s", a.Key, a.Value)
	}

	return fmt.Sprintf("%s=%c%s%c", a.Key, a.Quoted, a.Value, a.Quoted)
}
