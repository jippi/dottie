// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
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

func (s *File) RenderFull() string {
	return s.Render(RenderSettings{
		ShowPretty:       true,
		IncludeCommented: true,
	})
}

func (s *File) Render(config RenderSettings) string {
	return renderStatements(s.Statements, config)
}
