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

func (f *RenderSettings) Match(assignment *Assignment) bool {
	if len(f.FilterKeyPrefix) > 0 && !strings.HasPrefix(assignment.Key, f.FilterKeyPrefix) {
		return false
	}

	if !assignment.BelongsToGroup(*f) {
		return false
	}

	if assignment.Commented && !f.IncludeCommented {
		return false
	}

	return true
}

func (f *RenderSettings) WithComments() bool {
	return f.ShowPretty || f.ShowComments
}

func (f *RenderSettings) WithGroups() bool {
	return f.ShowPretty || f.ShowGroups
}

func (f *RenderSettings) WithBlankLines() bool {
	return f.ShowPretty || f.ShowBlankLines
}
