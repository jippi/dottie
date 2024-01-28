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

func (s *Assignment) Is(other Statement) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s *Assignment) BelongsToGroup(config RenderSettings) bool {
	return s.Group == nil || s.Group.BelongsToGroup(config)
}

func (s *Assignment) ShouldRender(config RenderSettings) bool {
	return config.Match(s) && s.BelongsToGroup(config)
}

func (s *Assignment) Render(config RenderSettings) string {
	var buff bytes.Buffer

	if config.WithComments() {
		buff.WriteString(s.String())
		buff.WriteString("\n")

		return buff.String()
	}

	buff.WriteString(s.Assignment())
	buff.WriteString("\n")

	return buff.String()
}

func (s *Assignment) statementNode() {}

func (s *Assignment) SetQuote(in string) {
	switch in {
	case "\"", "double":
		s.Quoted = DoubleQuotes
	case "'", "single":
		s.Quoted = SingleQuotes
	case "none":
		s.Quoted = NoQuotes
	}
}

func (s *Assignment) Assignment() string {
	if s.Quoted == 0 {
		return fmt.Sprintf("%s=%s", s.Key, s.Value)
	}

	return fmt.Sprintf("%s=%c%s%c", s.Key, s.Quoted, s.Value, s.Quoted)
}

func (s *Assignment) String() string {
	var buff bytes.Buffer

	for _, c := range s.Comments {
		buff.WriteString(c.Value)
		buff.WriteString("\n")
	}

	if s.Commented {
		buff.WriteString("#")
	}

	buff.WriteString(s.Assignment())

	return buff.String()
}
