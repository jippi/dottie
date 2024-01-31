package ast

import (
	"bytes"
	"strings"
	"unicode"
)

// Statement represents syntax tree node of .env file statement (like: assignment or comment).
type Statement interface {
	statementNode()
	Is(Statement) bool
	BelongsToGroup(RenderSettings) bool
	Render(RenderSettings) string
}

type Position struct {
	Line      uint
	FirstLine uint
	LastLine  uint
}

func renderStatements(statements []Statement, config RenderSettings) string {
	var buf bytes.Buffer
	var prev Statement
	var line *Newline

	var printed bool

	for _, stmt := range statements {
		switch val := stmt.(type) {

		case *Group:
			output := val.Render(config)
			if len(output) == 0 {
				continue
			}

			if config.WithBlankLines() && !prev.Is(line) {
				buf.WriteString("\n")
			}

			printed = true
			buf.WriteString(output)

		case *Comment:
			printed = true

			buf.WriteString(val.Render(config))

		case *Assignment:
			output := val.Render(config)
			if len(output) == 0 {
				continue
			}

			// Looks like current and previous is both "Assignment"
			// which mean they are too close in the document, so we will
			// attempt to inject some new-lines to give them some space
			if config.WithBlankLines() && val.Is(prev) {
				switch {

				// only allow cuddling of assignments if they both have no comments
				case val.HasComments() || assignmentHasComments(prev):
					buf.WriteString("\n")

				default:
					// NOOP
				}
			}

			printed = true
			buf.WriteString(output)

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
		str += "\n"
	}

	return str
}

func assignmentHasComments(stmt Statement) bool {
	x, ok := stmt.(*Assignment)
	if !ok {
		return false
	}

	return x.HasComments()
}
