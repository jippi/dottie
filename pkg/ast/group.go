package ast

import (
	"bytes"
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

func (g *Group) BelongsToGroup(config RenderSettings) bool {
	if len(config.FilterGroup) == 0 {
		return true
	}

	return g.String() == config.FilterGroup || slug.Make(g.String()) == config.FilterGroup
}

func (g *Group) String() string {
	return strings.TrimPrefix(g.Name, "# ")
}

func (g *Group) Render(config RenderSettings) string {
	if !g.BelongsToGroup(config) {
		return ""
	}

	var buff bytes.Buffer

	res := renderStatements(g.Statements, config)

	if config.WithGroups() && len(res) > 0 {
		buff.WriteString("################################################################################")
		buff.WriteString("\n")

		buff.WriteString(g.Name)
		buff.WriteString("\n")

		buff.WriteString("################################################################################")
		buff.WriteString("\n")
	}

	// Render the statements attached to the group
	buff.WriteString(res)

	return buff.String()
}
