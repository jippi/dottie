package ast

import (
	"reflect"
	"strings"

	"github.com/gosimple/slug"
)

type Group struct {
	Name       string      // Name of the group (within the header)
	Statements []Statement // Statements within the group
	Position   Position    // Positional information about the group
}

func (g *Group) statementNode() {
}

func (g *Group) Is(other Statement) bool {
	return reflect.TypeOf(g) == reflect.TypeOf(other)
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
