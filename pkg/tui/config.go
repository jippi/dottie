package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type ColorConfig struct {
	Text         ComponentColor
	TextEmphasis ComponentColor
	Background   ComponentColor
	Border       ComponentColor
}

func NewNeutralColorComponentConfig() ColorConfig {
	config := ColorConfig{
		Text:         ComponentColor{},
		TextEmphasis: ComponentColor{},
		Background:   ComponentColor{},
		Border:       ComponentColor{},
	}

	return config
}

func NewColorComponentConfig(baseColor lipgloss.Color) ColorConfig {
	base := ColorToHex(baseColor)

	config := ColorConfig{
		Text: ComponentColor{
			Light: ComponentColorConfig{
				Color: base,
			},
		},
		TextEmphasis: ComponentColor{
			Light: ComponentColorConfig{
				Color:   base,
				Filter:  "shade",
				Percent: 0.6,
			},
			Dark: ComponentColorConfig{
				Color:   base,
				Filter:  "tint",
				Percent: 0.4,
			},
		},
		Background: ComponentColor{
			Light: ComponentColorConfig{
				Color:   base,
				Filter:  "tint",
				Percent: 0.8,
			},
			Dark: ComponentColorConfig{
				Color:   base,
				Filter:  "shade",
				Percent: 0.8,
			},
		},
		Border: ComponentColor{
			Light: ComponentColorConfig{
				Color:   base,
				Filter:  "tint",
				Percent: 0.6,
			},
			Dark: ComponentColorConfig{
				Color:   base,
				Filter:  "shade",
				Percent: 0.4,
			},
		},
	}

	return config
}

type ComponentColor struct {
	Light ComponentColorConfig
	Dark  ComponentColorConfig
}

func (cc ComponentColor) AdaptiveColor() lipgloss.AdaptiveColor {
	result := lipgloss.AdaptiveColor{}
	result.Light = cc.Light.AsHex()
	result.Dark = cc.Dark.AsHex()

	return result
}

type ComponentColorConfig struct {
	Color    string
	MixColor string
	Filter   string
	Percent  float64
}

func (ccc ComponentColorConfig) AsHex() string {
	switch ccc.Filter {
	case "shade":
		return ColorToHex(ShadeColor(ccc.Color, ccc.Percent))

	case "tint":
		return ColorToHex(TintColor(ccc.Color, ccc.Percent))

	case "mix":
		percent := ccc.Percent
		if percent == 0 {
			percent = 0.5
		}

		return ColorToHex(MixColors(ccc.Color, ccc.MixColor, percent))

	default:
		return ccc.Color
	}
}
