package ast

import (
	"fmt"
)

// Statement represents syntax tree node of .env file statement (like: assignment or comment).
type Statement interface {
	statementNode()
	Is(statement Statement) bool
	Type() string
}

type Position struct {
	File      string `json:"file"`
	Line      uint   `json:"line"`
	FirstLine uint   `json:"first_line"`
	LastLine  uint   `json:"last_line"`
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d", p.File, p.Line)
}
