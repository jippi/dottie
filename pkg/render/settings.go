package render

type OutputType uint

const (
	Plain OutputType = iota
	Colorized
	CompletionKeyOnly
)

type Settings struct {
	retainKeyPrefix    string
	retainGroup        string
	includeDisabled    bool
	showBlankLines     bool
	showColors         bool
	showComments       bool
	ShowGroupBanners   bool
	formatOutput       bool
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

func (rs Settings) Handlers() []Handler {
	var res []Handler

	if !rs.showComments {
		res = append(res, ExcludeComments)
	}

	if !rs.includeDisabled {
		res = append(res, ExcludeDisabledAssignments)
	}

	if len(rs.retainGroup) > 0 {
		res = append(res, RetainGroup(rs.retainGroup))
	}

	if len(rs.retainKeyPrefix) > 0 {
		res = append(res, RetainKeyPrefix(rs.retainKeyPrefix))
	}

	return res
}
