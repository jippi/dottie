package render

type Settings struct {
	filterKeyPrefix       string
	filterGroup           string
	includeDisabled       bool
	showBlankLines        bool
	showColors            bool
	showComments          bool
	ShowGroupBanners      bool
	formatOutput          bool
	useInterpolatedValues bool
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
