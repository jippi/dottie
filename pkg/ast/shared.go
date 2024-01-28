package ast

import (
	"bytes"
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

	for _, stmt := range statements {
		switch val := stmt.(type) {
		case *Group:
			if !val.ShouldRender(config) {
				continue
			}

			buff.WriteString(val.Render(config))

			previous = stmt

		case *Comment:
			if !val.ShouldRender(config) {
				continue
			}

			previous = stmt

			buff.WriteString(val.Render(config))

		case *Assignment:
			if !val.ShouldRender(config) {
				continue
			}

			previous = stmt
			buff.WriteString(val.Render(config))

		case *Newline:
			if !val.ShouldRender(config) {
				continue
			}

			// Don't print multiple newlines after each other
			if val.Is(previous) {
				continue
			}

			previous = stmt

			buff.WriteString(val.Render(config))
		}
	}

	return buff.String()
}
