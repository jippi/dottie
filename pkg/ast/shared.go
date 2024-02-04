package ast

import (
	"fmt"
)

// Statement represents syntax tree node of .env file statement (like: assignment or comment).
type Statement interface {
	statementNode()
	Is(Statement) bool
}

type Position struct {
	File      string
	Line      uint
	FirstLine uint
	LastLine  uint
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d", p.File, p.Line)
}
