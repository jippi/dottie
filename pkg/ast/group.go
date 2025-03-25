package ast

import (
	"reflect"
	"strings"

	"github.com/gosimple/slug"
	"github.com/ryanuber/go-glob"
)

type Group struct {
	Name       string      `json:"name"`       // Name of the group (within the header)
	Statements []Statement `json:"statements"` // Statements within the group
	Position   Position    `json:"position"`   // Positional information about the group
}

func (g *Group) statementNode() {
}

func (g *Group) Is(other Statement) bool {
	if g == nil || other == nil {
		return false
	}

	return g.Type() == other.Type()
}

func (g *Group) Type() string {
	if g == nil {
		return "<nil>Group"
	}

	return reflect.TypeOf(g).String()
}

func (g *Group) BelongsToGroup(name string) bool {
	if len(name) == 0 {
		return true
	}

	return glob.Glob(name, g.String()) || glob.Glob(name, slug.Make(g.String()))
}

func (g *Group) String() string {
	return strings.TrimPrefix(g.Name, "# ")
}

func (g *Group) Assignments() []*Assignment {
	var assignments []*Assignment

	for _, statement := range g.Statements {
		if assign, ok := statement.(*Assignment); ok {
			assignments = append(assignments, assign)
		}
	}

	return assignments
}

func (g *Group) GetAssignmentIndex(name string) (int, *Assignment) {
	for i, assign := range g.Assignments() {
		if assign.Name == name {
			return i, assign
		}
	}

	return -1, nil
}
