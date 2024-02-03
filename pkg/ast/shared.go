package ast

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// Statement represents syntax tree node of .env file statement (like: assignment or comment).
type Statement interface {
	statementNode()
	Is(Statement) bool
	BelongsToGroup(RenderSettings) bool
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

func renderStatements(statements []Statement, config RenderSettings) string {
	var (
		buf     bytes.Buffer
		prev    Statement
		printed bool
	)

	for _, stmt := range statements {
		switch val := stmt.(type) {
		case *Group:
			panic("group should never happen in renderStatements")

		case *Comment:
			printed = true

			buf.WriteString(val.Render(config, false))

		case *Assignment:
			output := val.Render(config)
			if len(output) == 0 {
				continue
			}

			// Looks like current and previous is both "Assignment"
			// which mean they are too close in the document, so we will
			// attempt to inject some new-lines to give them some space
			if config.WithBlankLines() && val.Is(prev) {
				// only allow cuddling of assignments if they both have no comments
				if val.HasComments() || assignmentHasComments(prev) {
					buf.WriteString("\n")
				}
			}

			buf.WriteString(output)

			printed = true

		case *Newline:
			output := val.Render(config)
			if len(output) == 0 {
				continue
			}

			// Don't print multiple newlines after each other
			if val.Is(prev) {
				continue
			}

			buf.WriteString(output)
		}

		prev = stmt
	}

	// If nothing "useful" was printed, don't bother outputting the groups buffer
	if !printed {
		return ""
	}

	str := buf.String()

	// Remove any duplicate newlines that might have crept into the output
	if config.WithBlankLines() {
		str = strings.TrimRightFunc(str, unicode.IsSpace)
	}

	return "\n" + str
}

func assignmentHasComments(stmt Statement) bool {
	x, ok := stmt.(*Assignment)
	if !ok {
		return false
	}

	return x.HasComments()
}
