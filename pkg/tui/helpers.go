package tui

import (
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/teacat/noire"
)

func Renderer(w io.Writer, opts ...termenv.OutputOption) *lipgloss.Renderer {
	return lipgloss.NewRenderer(w, opts...)
}

func RendererWithTTY(w io.Writer, opts ...termenv.OutputOption) *lipgloss.Renderer {
	opts = append(opts, termenv.WithTTY(true))

	return lipgloss.NewRenderer(w, opts...)
}

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

func MixColors(a, b string, weight float64) lipgloss.Color {
	if weight < 0 || weight > 1 {
		panic("MixColors [weight] must be between 0.0 and 1.0 (0.5 == 50%)")
	}

	return lipgloss.Color("#" + noire.NewHex(a).Mix(noire.NewHex(b), weight).Hex())
}

func ColorToHex(in lipgloss.Color) string {
	return string(in)
}
