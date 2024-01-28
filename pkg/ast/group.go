package ast

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/gosimple/slug"
)

type Group struct {
	Name       string
	FirstLine  int
	LastLine   int
	Statements []Statement
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

func (g *Group) statementNode() {
}

func (g *Group) String() string {
	return strings.TrimPrefix(g.Name, "# ")
}

func (g *Group) ShouldRender(config RenderSettings) bool {
	if !g.BelongsToGroup(config) {
		return false
	}

	for _, stmt := range g.Statements {
		switch val := stmt.(type) {
		case *Assignment:
			if !val.ShouldRender(config) {
				continue
			}

			if config.Match(val) {
				return true
			}

		case *Comment:
			if val.ShouldRender(config) {
				return true
			}
		}
	}

	return false
}

func (g *Group) Render(config RenderSettings) string {
	var buff bytes.Buffer

	if config.WithGroups() {
		buff.WriteString("################################################################################")
		buff.WriteString("\n")

		buff.WriteString(g.Name)
		buff.WriteString("\n")

		buff.WriteString("################################################################################")
		buff.WriteString("\n")
	}

	buff.WriteString(renderStatements(g.Statements, config))

	return buff.String()
}
