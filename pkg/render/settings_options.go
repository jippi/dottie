package render

type SettingsOption func(*Settings)

func WithFilterKeyPrefix(prefix string) SettingsOption {
	return func(s *Settings) {
		s.retainKeyPrefix = prefix
	}
}

func WithFilterGroup(name string) SettingsOption {
	return func(s *Settings) {
		s.retainGroup = name
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
		s.showBlankLines = boolean
		s.showColors = boolean
		s.showComments = boolean
		s.ShowGroupBanners = boolean
	}
}

func WithExport(boolean bool) SettingsOption {
	return func(s *Settings) {
		s.export = boolean
	}
}

func WithInterpolation(b bool) SettingsOption {
	return func(s *Settings) {
		s.InterpolatedValues = b
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
