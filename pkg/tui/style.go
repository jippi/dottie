package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type Style struct {
	Text         lipgloss.AdaptiveColor
	TextEmphasis lipgloss.AdaptiveColor
	Background   lipgloss.AdaptiveColor
	Border       lipgloss.AdaptiveColor

	noColor bool
}

func NewStyle(baseColor lipgloss.Color) Style {
	base := ColorToHex(baseColor)

	style := Style{
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

	if len(style.Text.Dark) == 0 {
		style.Text.Dark = style.TextEmphasis.Dark
	}

	return style
}

func NewStyleWithoutColor() Style {
	return Style{
		noColor: true,
	}
}

func (style Style) NewPrinter(renderer *lipgloss.Renderer, options ...PrinterOption) Print {
	return NewPrinter(style, renderer, options...)
}

func (style Style) TextStyle() lipgloss.Style {
	if style.noColor {
		return lipgloss.NewStyle()
	}

	return lipgloss.
		NewStyle().
		Foreground(style.Text)
}

func (c Style) TextEmphasisStyle() lipgloss.Style {
	if c.noColor {
		return lipgloss.NewStyle()
	}

	return lipgloss.NewStyle().
		Foreground(c.TextEmphasis).
		Background(c.Background).
		Bold(true).
		BorderForeground(c.Border)
}

func (c Style) BoxStyles(header, body lipgloss.Style) Box {
	return Box{
		Header: header.
			Align(lipgloss.Center, lipgloss.Center).
			Border(headerBorder).
			BorderForeground(c.Border).
			PaddingBottom(1).
			PaddingTop(1).
			Inherit(c.TextEmphasisStyle()),

		Body: body.
			Align(lipgloss.Left).
			Border(bodyBorder).
			BorderForeground(c.Border).
			BorderTop(false).
			Padding(1),
	}
}
