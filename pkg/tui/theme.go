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

type ThemeConfig struct {
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

func (tc ThemeConfig) Printer(w io.Writer) ThemePrinter {
	return ThemePrinter{
		w:     w,
		theme: tc,
		cache: make(map[colorType]Printer),
	}
}

type ThemePrinter struct {
	theme ThemeConfig
	w     io.Writer
	cache map[colorType]Printer
}

func (tp ThemePrinter) Color(colorType colorType) Printer {
	if printer, ok := tp.cache[colorType]; ok {
		return printer
	}

	var color Color

	switch colorType {
	case Danger:
		color = tp.theme.Danger
	case Dark:
		color = tp.theme.Danger
	case Info:
		color = tp.theme.Info
	case Light:
		color = tp.theme.Light
	case Primary:
		color = tp.theme.Primary
	case Secondary:
		color = tp.theme.Secondary
	case Success:
		color = tp.theme.Success
	case Warning:
		color = tp.theme.Warning
	case NoColor:
		color = tp.theme.NoColor
	}

	tp.cache[colorType] = color.BufferPrinter(tp.w)

	return tp.cache[colorType]
}

func NewTheme() ThemeConfig {
	theme := ThemeConfig{}
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
