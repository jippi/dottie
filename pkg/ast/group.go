package ast

import (
	"reflect"
	"strings"

	"github.com/gosimple/slug"
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

	return g.String() == name || slug.Make(g.String()) == name
}

func (g *Group) String() string {
	return strings.TrimPrefix(g.Name, "# ")
}
