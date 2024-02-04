package render

import (
	"strings"

	"github.com/jippi/dottie/pkg/ast"
)

type Settings struct {
	FilterKeyPrefix  string
	FilterGroup      string
	IncludeCommented bool

	ShowBlankLines bool
	ShowColors     bool
	ShowComments   bool
	ShowGroups     bool
	ShowPretty     bool

	Interpolate bool
}

func (rs *Settings) Match(assignment *ast.Assignment) bool {
	if !assignment.Active && !rs.IncludeCommented {
		return false
	}

	if len(rs.FilterKeyPrefix) > 0 && !strings.HasPrefix(assignment.Name, rs.FilterKeyPrefix) {
		return false
	}

	if !assignment.BelongsToGroup(rs.FilterGroup) {
		return false
	}

	return true
}

func (rs *Settings) WithComments() bool {
	return rs.ShowPretty || rs.ShowComments
}

func (rs *Settings) WithGroups() bool {
	return rs.ShowPretty || rs.ShowGroups
}

func (rs *Settings) WithBlankLines() bool {
	return rs.ShowPretty || (rs.ShowBlankLines && rs.ShowComments)
}

func (rs *Settings) WithColors() bool {
	return rs.ShowColors
}
