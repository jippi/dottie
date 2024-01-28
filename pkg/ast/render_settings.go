package ast

import (
	"strings"
)

type RenderSettings struct {
	FilterKeyPrefix string
	FilterGroup     string
	FilterCommented bool

	ShowPretty     bool
	ShowComments   bool
	ShowGroups     bool
	ShowBlankLines bool
}

func (f *RenderSettings) Match(assignment *Assignment) bool {
	if len(f.FilterKeyPrefix) > 0 && !strings.HasPrefix(assignment.Key, f.FilterKeyPrefix) {
		return false
	}

	if len(f.FilterGroup) > 0 && assignment.Group != nil && assignment.Group.Name != f.FilterGroup {
		return false
	}

	if assignment.Commented && !f.FilterCommented {
		return false
	}

	return true
}

func (f *RenderSettings) Comments() bool {
	return f.ShowPretty || f.ShowComments
}

func (f *RenderSettings) Groups() bool {
	return f.ShowPretty || f.ShowGroups
}

func (f *RenderSettings) BlankLines() bool {
	return f.ShowPretty || f.ShowBlankLines
}
