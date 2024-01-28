package ast

import (
	"strings"
)

type RenderSettings struct {
	FilterKeyPrefix  string
	FilterGroup      string
	IncludeCommented bool

	ShowPretty     bool
	ShowComments   bool
	ShowGroups     bool
	ShowBlankLines bool
}

func (rs *RenderSettings) Match(assignment *Assignment) bool {
	if len(rs.FilterKeyPrefix) > 0 && !strings.HasPrefix(assignment.Key, rs.FilterKeyPrefix) {
		return false
	}

	if !assignment.BelongsToGroup(*rs) {
		return false
	}

	if assignment.Commented && !rs.IncludeCommented {
		return false
	}

	return true
}

func (rs *RenderSettings) WithComments() bool {
	return rs.ShowPretty || rs.ShowComments
}

func (rs *RenderSettings) WithGroups() bool {
	return rs.ShowPretty || rs.ShowGroups
}

func (rs *RenderSettings) WithBlankLines() bool {
	return rs.ShowPretty || rs.ShowBlankLines
}
