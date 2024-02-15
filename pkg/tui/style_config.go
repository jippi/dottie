package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type StyleConfig struct {
	Text         lipgloss.AdaptiveColor
	TextEmphasis lipgloss.AdaptiveColor
	Background   lipgloss.AdaptiveColor
	Border       lipgloss.AdaptiveColor
}

func NewStyleConfig(baseColor lipgloss.Color) StyleConfig {
	base := ColorToHex(baseColor)

	config := StyleConfig{
		Text: lipgloss.AdaptiveColor{
			Light: transformColor(base, "", 0),
		},
		TextEmphasis: lipgloss.AdaptiveColor{
			Light: transformColor(base, "shade", 0.6),
			Dark:  transformColor(base, "tint", 0.4),
		},
		Background: lipgloss.AdaptiveColor{
			Light: transformColor(base, "tint", 0.8),
			Dark:  transformColor(base, "shade", 0.8),
		},
		Border: lipgloss.AdaptiveColor{
			Light: transformColor(base, "tint", 0.6),
			Dark:  transformColor(base, "shade", 0.4),
		},
	}

	return config
}

func transformColor(base, filter string, percent float64) string {
	switch filter {
	case "shade":
		return ColorToHex(ShadeColor(base, percent))

	case "tint":
		return ColorToHex(TintColor(base, percent))

	case "mix":
		panic("unexpected mix filter")

	default:
		return base
	}
}
