package tui

import (
	"io"
)

type colorType int

const (
	Danger colorType = iota
	Dark
	Info
	Light
	NoColor
	Primary
	Secondary
	Success
	Warning
)

type Theme struct {
	DefaultWidth int

	Danger    Color
	Dark      Color
	Info      Color
	Light     Color
	NoColor   Color
	Primary   Color
	Secondary Color
	Success   Color
	Warning   Color
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
	theme.DefaultWidth = 80

	theme.Danger = NewColor(NewColorComponentConfig(Red))
	theme.Info = NewColor(NewColorComponentConfig(Cyan))
	theme.Light = NewColor(NewColorComponentConfig(Gray300))
	theme.Primary = NewColor(NewColorComponentConfig(Blue))
	theme.Secondary = NewColor(NewColorComponentConfig(Gray600))
	theme.Success = NewColor(NewColorComponentConfig(Green))
	theme.Warning = NewColor(NewColorComponentConfig(Yellow))
	theme.NoColor = NewNoColor()

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
