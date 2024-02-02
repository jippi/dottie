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

	Interpolate bool
}

func (rs *RenderSettings) Match(assignment *Assignment) bool {
	if !assignment.Active && !rs.IncludeCommented {
		return false
	}

	if len(rs.FilterKeyPrefix) > 0 && !strings.HasPrefix(assignment.Name, rs.FilterKeyPrefix) {
		return false
	}

	if !assignment.BelongsToGroup(*rs) {
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
	return rs.ShowPretty || (rs.ShowBlankLines && rs.ShowComments)
}
