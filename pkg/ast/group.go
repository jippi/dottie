package ast

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/gosimple/slug"
	"github.com/jippi/dottie/pkg/tui"
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

	var buf bytes.Buffer

	res := renderStatements(g.Statements, config)

	if config.WithGroups() && len(res) > 0 {
		if config.WithColors() {
			out := tui.Theme.Info.Printer(tui.RendererWithTTY(&buf))
			out.Println("################################################################################")
			out.ApplyStyle(tui.Bold).Println(g.Name)
			out.Println("################################################################################")
		} else {
			buf.WriteString("################################################################################")
			buf.WriteString("\n")

			buf.WriteString(g.Name)
			buf.WriteString("\n")

			buf.WriteString("################################################################################")
			buf.WriteString("\n")
		}
	}

	// Render the statements attached to the group
	buf.WriteString(res)

	return buf.String()
}
