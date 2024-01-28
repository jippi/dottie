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
	Key       string
	Value     string
	Comments  []*Comment
	Group     *Group
	Commented bool

	FirstLine  int
	LastLine   int
	LineNumber int
	Naked      bool
	Complete   bool
	Quoted     rune
}

func (a *Assignment) Is(other Statement) bool {
	if other == nil {
		return false
	}

	return reflect.TypeOf(a) == reflect.TypeOf(other)
}

func (a *Assignment) BelongsToGroup(config RenderSettings) bool {
	return a.Group == nil || a.Group.BelongsToGroup(config)
}

func (a *Assignment) HasComment() bool {
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

func (a *Assignment) statementNode() {}

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
	if a.Quoted == 0 {
		return fmt.Sprintf("%s=%s", a.Key, a.Value)
	}

	return fmt.Sprintf("%s=%c%s%c", a.Key, a.Quoted, a.Value, a.Quoted)
}
