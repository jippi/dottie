package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/teacat/noire"
)

func ShadeColor(in string, percent float64) lipgloss.Color {
	if percent < 0 || percent > 1 {
		panic("ShadeColor [percent] must be between 0.0 and 1.0 (0.5 == 50%)")
	}

	return lipgloss.Color("#" + noire.NewHex(in).Shade(percent).Hex())
}

func TintColor(in string, percent float64) lipgloss.Color {
	if percent < 0 || percent > 1 {
		panic("TintColor [percent] must be between 0.0 and 1.0 (0.5 == 50%)")
	}

	return lipgloss.Color("#" + noire.NewHex(in).Tint(percent).Hex())
}

func ColorToHex(in lipgloss.Color) string {
	return string(in)
}

func TransformColor(base, filter string, percent float64) string {
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
