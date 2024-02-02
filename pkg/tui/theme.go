package tui

import (
	"github.com/erikgeiser/promptkit"
)

type ThemeConfig struct {
	DefaultWidth int

	// Line wrapping handling
	WrapMode promptkit.WrapMode

	Danger    Color
	Dark      Color
	Info      Color
	Light     Color
	Primary   Color
	Secondary Color
	Success   Color
	Warning   Color
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
