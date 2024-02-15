package tui

import (
	"io"

	"github.com/erikgeiser/promptkit"
)

type colorType int

const (
	Danger colorType = iota
	Dark
	Info
	Light
	Neutral
	Primary
	Secondary
	Success
	Warning
)

type ThemeConfig struct {
	DefaultWidth int

	// Line wrapping handling
	WrapMode promptkit.WrapMode

	Danger    Color
	Dark      Color
	Info      Color
	Light     Color
	Neutral   Color
	Primary   Color
	Secondary Color
	Success   Color
	Warning   Color
}

func (tc ThemeConfig) Printer(w io.Writer) ThemePrinter {
	return ThemePrinter{
		w:     w,
		cache: make(map[colorType]Printer),
	}
}

type ThemePrinter struct {
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
		color = Theme.Danger
	case Dark:
		color = Theme.Danger
	case Info:
		color = Theme.Info
	case Light:
		color = Theme.Light
	case Primary:
		color = Theme.Primary
	case Secondary:
		color = Theme.Secondary
	case Success:
		color = Theme.Success
	case Warning:
		color = Theme.Warning
	case Neutral:
		color = Theme.Neutral
	}

	tp.cache[colorType] = color.BuffPrinter(tp.w)

	return tp.cache[colorType]
}

var Theme ThemeConfig

func init() {
	Theme = ThemeConfig{}
	Theme.DefaultWidth = 100
	Theme.WrapMode = nil // Disabled for now, left here for easy opt-in in the future

	Theme.Danger = NewColor(NewColorComponentConfig(Red))
	Theme.Info = NewColor(NewColorComponentConfig(Cyan))
	Theme.Light = NewColor(NewColorComponentConfig(Gray300))
	Theme.Primary = NewColor(NewColorComponentConfig(Blue))
	Theme.Secondary = NewColor(NewColorComponentConfig(Gray600))
	Theme.Success = NewColor(NewColorComponentConfig(Green))
	Theme.Warning = NewColor(NewColorComponentConfig(Yellow))
	Theme.Neutral = NewColor(NewNeutralColorComponentConfig())

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

	Theme.Dark = NewColor(dark)
}
