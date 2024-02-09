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
		if b {
			s.outputter = ColorizedOutput{}
		} else {
			s.outputter = PlainOutput{}
		}
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

func WithFormattedOutput(boolean bool) SettingsOption {
	return func(s *Settings) {
		s.formatOutput = boolean
		s.showComments = boolean
		s.ShowGroupBanners = boolean
		s.showColors = boolean
		s.showBlankLines = boolean
	}
}

func WithInterpolation(b bool) SettingsOption {
	return func(s *Settings) {
		s.useInterpolatedValues = b
	}
}

func WithOutputter(o Output) SettingsOption {
	return func(s *Settings) {
		s.outputter = o
	}
}

func WithOutputType(t OutputType) SettingsOption {
	return func(settings *Settings) {
		switch t {
		case Plain:
			settings.outputter = PlainOutput{}

		case Colorized:
			settings.outputter = ColorizedOutput{}

		case CompletionKeyOnly:
			settings.outputter = CompletionOutputKeys{}

		default:
			panic("Invalid outputter type")
		}
	}
}
