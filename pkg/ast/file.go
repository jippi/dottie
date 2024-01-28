// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"bytes"
	"reflect"
)

// File node represents .env file statement, that contains assignments and comments.
type File struct {
	Statements []Statement
	Groups     []*Group
}

func (s *File) Is(other Statement) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s *File) BelongsToGroup(config RenderSettings) bool {
	return false
}

func (s *File) statementNode() {
}

func (s *File) Pairs() map[string]string {
	values := map[string]string{}

	for _, stmt := range s.Statements {
		if assign, ok := stmt.(*Assignment); ok {
			values[assign.Key] = assign.Value
		}
	}

	return values
}

func (s *File) GetGroup(config RenderSettings) *Group {
	for _, grp := range s.Groups {
		if grp.BelongsToGroup(config) {
			return grp
		}
	}

	return nil
}

func (s *File) Get(name string) *Assignment {
	for _, stmt := range s.Statements {
		if assign, ok := stmt.(*Assignment); ok {
			if assign.Key == name {
				return assign
			}
		}
	}

	return nil
}

func (s *File) ShouldRender(config RenderSettings) bool {
	return true
}

func (s *File) Render() []byte {
	return s.RenderWithFilter(RenderSettings{
		ShowPretty:       true,
		IncludeCommented: true,
	})
}

func (s *File) RenderWithFilter(config RenderSettings) []byte {
	var buff bytes.Buffer
	var previous Statement

	for _, stmt := range s.Statements {
		switch val := stmt.(type) {
		case *Group:
			if !val.ShouldRender(config) {
				continue
			}

			previous = stmt

			buff.WriteString("################################################################################")
			buff.WriteString("\n")

			buff.WriteString(val.Name)
			buff.WriteString("\n")

			buff.WriteString("################################################################################")
			buff.WriteString("\n")

		case *Comment:
			if !val.ShouldRender(config) {
				continue
			}

			previous = stmt

			buff.WriteString(val.String())
			buff.WriteString("\n")

		case *Assignment:
			if !val.ShouldRender(config) {
				continue
			}

			previous = stmt

			if config.WithComments() {
				buff.WriteString(val.String())
				buff.WriteString("\n")
				// buff.WriteString("\n")

				continue
			}

			buff.WriteString(val.Assignment())
			buff.WriteString("\n")

		case *Newline:
			if !val.ShouldRender(config) {
				continue
			}

			// Don't print multiple newlines after each other
			if val.Is(previous) {
				continue
			}

			previous = stmt

			buff.WriteString("\n")
		}
	}

	b := bytes.TrimSpace(buff.Bytes())
	b = append(b, byte('\n'))

	return b
}
