package tui

import (
	"io"
)

type Theme struct {
	Danger    Style
	Dark      Style
	Info      Style
	Light     Style
	NoColor   Style
	Primary   Style
	Secondary Style
	Success   Style
	Warning   Style
}

func (theme Theme) Printer(w io.Writer) ThemeWriter {
	return ThemeWriter{
		w:     w,
		theme: theme,
		cache: make(map[colorType]Printer),
	}
}

func NewTheme() Theme {
	theme := Theme{}

	theme.Danger = NewColor(NewColorComponentConfig(Red))
	theme.Info = NewColor(NewColorComponentConfig(Cyan))
	theme.Light = NewColor(NewColorComponentConfig(Gray300))
	theme.Primary = NewColor(NewColorComponentConfig(Blue))
	theme.Secondary = NewColor(NewColorComponentConfig(Gray600))
	theme.Success = NewColor(NewColorComponentConfig(Green))
	theme.Warning = NewColor(NewColorComponentConfig(Yellow))
	theme.NoColor = NewStyle()

	dark := NewColorComponentConfig(Gray700)

	dark.TextEmphasis.Dark = ComponentColorConfig{
		Color: ColorToHex(Gray300),
	}

	dark.Background.Dark = ComponentColorConfig{
		Color: "#1a1d20",
	}

	dark.Border.Dark = ComponentColorConfig{
		Color: ColorToHex(Gray800),
	}

	theme.Dark = NewColor(dark)

	return theme
}
