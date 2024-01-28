// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"reflect"
)

// Document node represents .env file statement, that contains assignments and comments.
type Document struct {
	Statements  []Statement
	Groups      []*Group
	Assignments []*Assignment
}

func (d *Document) Is(other Statement) bool {
	return reflect.TypeOf(d) == reflect.TypeOf(other)
}

func (d *Document) BelongsToGroup(config RenderSettings) bool {
	return false
}

func (d *Document) statementNode() {
}

func (d *Document) Pairs() map[string]string {
	values := map[string]string{}

	for _, stmt := range d.Statements {
		if assign, ok := stmt.(*Assignment); ok {
			values[assign.Key] = assign.Value
		}
	}

	return values
}

func (d *Document) GetGroup(config RenderSettings) *Group {
	for _, grp := range d.Groups {
		if grp.BelongsToGroup(config) {
			return grp
		}
	}

	return nil
}

func (d *Document) Get(name string) *Assignment {
	for _, assign := range d.Assignments {
		if assign.Key == name {
			return assign
		}
	}

	return nil
}

func (d *Document) GetPosition(name string) (int, *Assignment) {
	for i, assign := range d.Assignments {
		if assign.Key == name {
			return i, assign
		}
	}

	return -1, nil
}

func (d *Document) ShouldRender(config RenderSettings) bool {
	return true
}

func (d *Document) RenderFull() string {
	return d.Render(RenderSettings{
		ShowPretty:       true,
		IncludeCommented: true,
	})
}

func (d *Document) Render(config RenderSettings) string {
	return renderStatements(d.Statements, config)
}
