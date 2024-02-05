package render

type SettingsOption func(*Settings)

func WithFilterKeyPrefix(prefix string) SettingsOption {
	return func(s *Settings) {
		s.filterKeyPrefix = prefix
	}
}

func WithFilterGroup(name string) SettingsOption {
	return func(s *Settings) {
		s.filterGroup = name
	}
}

func WithIncludeDisabled(b bool) SettingsOption {
	return func(s *Settings) {
		s.includeDisabled = b
	}
}

func WithBlankLines(b bool) SettingsOption {
	return func(s *Settings) {
		s.showBlankLines = b
	}
}

func WithColors(b bool) SettingsOption {
	return func(s *Settings) {
		s.showColors = b
	}
}

func WithComments(b bool) SettingsOption {
	return func(s *Settings) {
		s.showComments = b
	}
}

func WithGroupBanners(b bool) SettingsOption {
	return func(s *Settings) {
		s.ShowGroupBanners = b
	}
}

func WithFormattedOutput(b bool) SettingsOption {
	return func(s *Settings) {
		s.formatOutput = b
		s.showComments = b
		s.ShowGroupBanners = b
		s.showColors = b
		s.showBlankLines = b
	}
}

func WithInterpolation(b bool) SettingsOption {
	return func(s *Settings) {
		s.useInterpolatedValues = b
	}
}
