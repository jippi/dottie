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
	ShouldRender(RenderSettings) bool
	Render(RenderSettings) string
}

func renderStatements(statements []Statement, config RenderSettings) string {
	var buff bytes.Buffer
	var previous Statement
	var line *Newline

	for _, stmt := range statements {
		switch val := stmt.(type) {
		case *Group:
			if !val.ShouldRender(config) {
				continue
			}

			if config.WithBlankLines() && !previous.Is(line) {
				buff.WriteString("\n")
			}

			buff.WriteString(val.Render(config))

		case *Comment:
			if !val.ShouldRender(config) {
				continue
			}

			buff.WriteString(val.Render(config))

		case *Assignment:
			if !val.ShouldRender(config) {
				continue
			}

			// Avoid assignments with comments cuddling
			if config.WithBlankLines() && val.Is(previous) {
				switch {
				// only allow cuddling of assignments if they both have no comments
				case !val.HasComment() && !hasComment(previous):

					// otherwise add some spacing
				default:
					buff.WriteString("\n")
				}
			}

			buff.WriteString(val.Render(config))

		case *Newline:
			if !val.ShouldRender(config) {
				continue
			}

			// Don't print multiple newlines after each other
			if val.Is(previous) {
				continue
			}

			buff.WriteString(val.Render(config))
		}

		previous = stmt
	}

	str := buff.String()

	if config.WithBlankLines() {
		str = strings.TrimRightFunc(str, unicode.IsSpace)
		str += "\n"
	}

	return str
}

func hasComment(stmt Statement) bool {
	x, ok := stmt.(*Assignment)
	if !ok {
		return false
	}

	return x.HasComment()
}
