package render

type Settings struct {
	FilterKeyPrefix       string
	FilterGroup           string
	IncludeDisabled       bool
	ShowBlankLines        bool
	ShowColors            bool
	ShowComments          bool
	ShowGroupBanners      bool
	FormatOutput          bool
	UseInterpolatedValues bool
}

type SettingsOption func(*Settings)

func NewSettings(options ...SettingsOption) *Settings {
	settings := &Settings{}

	for _, option := range options {
		option(settings)
	}

	return settings
}

func WithFilterKeyPrefix(prefix string) SettingsOption {
	return func(s *Settings) {
		s.FilterKeyPrefix = prefix
	}
}

func WithFilterGroup(name string) SettingsOption {
	return func(s *Settings) {
		s.FilterGroup = name
	}
}

func WithIncludeDisabled(b bool) SettingsOption {
	return func(s *Settings) {
		s.IncludeDisabled = b
	}
}

func WithBlankLines(b bool) SettingsOption {
	return func(s *Settings) {
		s.ShowBlankLines = b
	}
}

func WithColors(b bool) SettingsOption {
	return func(s *Settings) {
		s.ShowColors = b
	}
}

func WithComments(b bool) SettingsOption {
	return func(s *Settings) {
		s.ShowComments = b
	}
}

func WithGroupBanners(b bool) SettingsOption {
	return func(s *Settings) {
		s.ShowGroupBanners = b
	}
}

func WithFormattedOutput(b bool) SettingsOption {
	return func(s *Settings) {
		s.FormatOutput = b
		s.ShowComments = b
		s.ShowGroupBanners = b
		s.ShowColors = b
	}
}

func WithInterpolation(b bool) SettingsOption {
	return func(s *Settings) {
		s.UseInterpolatedValues = b
	}
}

func (rs Settings) WithBlankLines() bool {
	return !rs.FormatOutput && (rs.ShowBlankLines && rs.ShowComments)
}
