package render

import "github.com/jippi/dottie/pkg/ast"

type OutputType uint

const (
	Plain OutputType = iota
	Colorized
	CompletionKeyOnly
)

type Settings struct {
	retainKeys         []string
	retainKeyPrefix    string
	retainGroup        string
	includeDisabled    bool
	showBlankLines     bool
	showColors         bool
	showComments       bool
	ShowGroupBanners   bool
	formatOutput       bool
	export             bool
	InterpolatedValues bool
	outputter          Output
}

func NewSettings(options ...SettingsOption) *Settings {
	settings := &Settings{}

	return settings.Apply(options...)
}

func (s *Settings) Apply(options ...SettingsOption) *Settings {
	for _, option := range options {
		option(s)
	}

	return s
}

func (rs Settings) ShowBlankLines() bool {
	return rs.formatOutput || (rs.showBlankLines && rs.showComments)
}

func (rs Settings) Handlers() []ast.Selector {
	var res []ast.Selector

	if !rs.showComments {
		res = append(res, ast.ExcludeComments)
	}

	if !rs.includeDisabled {
		res = append(res, ast.ExcludeDisabledAssignments)
	}

	if len(rs.retainGroup) > 0 {
		res = append(res, ast.RetainGroup(rs.retainGroup))
	}

	if len(rs.retainKeyPrefix) > 0 {
		res = append(res, ast.RetainKeyPrefix(rs.retainKeyPrefix))
	}

	if len(rs.retainKeys) > 0 {
		res = append(res, ast.RetainExactKey(rs.retainKeys...))
	}

	return res
}
