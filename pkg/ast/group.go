package ast

import (
	"reflect"
	"strings"
)

type Group struct {
	Name       string
	FirstLine  int
	LastLine   int
	Statements []Statement
}

func (s *Group) Is(other Statement) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s *Group) BelongsToGroup(config RenderSettings) bool {
	if len(config.FilterGroup) == 0 {
		return true
	}

	return s.String() == config.FilterGroup
}

func (s *Group) statementNode() {
}

func (s *Group) String() string {
	return strings.TrimPrefix(s.Name, "# ")
}

func (s *Group) ShouldRender(config RenderSettings) bool {
	if !config.WithGroups() || !s.BelongsToGroup(config) {
		return false
	}

	for _, stmt := range s.Statements {
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
