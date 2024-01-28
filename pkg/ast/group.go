package ast

import (
	"strings"
)

type Group struct {
	Name       string
	FirstLine  int
	LastLine   int
	Statements []Statement
}

func (s *Group) statementNode() {
}

func (s *Group) String() string {
	return strings.TrimPrefix(s.Name, "# ")
}
