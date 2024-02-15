package tui

import (
	"io"

	"github.com/charmbracelet/lipgloss"
)

type Color struct {
	Text         lipgloss.AdaptiveColor
	TextEmphasis lipgloss.AdaptiveColor
	Background   lipgloss.AdaptiveColor
	Border       lipgloss.AdaptiveColor

	noColor bool
}

func NewColor(config ColorConfig) Color {
	color := Color{
		Text:         config.Text.AdaptiveColor(),
		TextEmphasis: config.TextEmphasis.AdaptiveColor(),
		Background:   config.Background.AdaptiveColor(),
		Border:       config.Border.AdaptiveColor(),
	}

	if len(color.Text.Dark) == 0 {
		color.Text.Dark = color.TextEmphasis.Dark
	}

	return color
}

func NewNoColor() Color {
	return Color{
		noColor: true,
	}
}

func (c Color) Printer(renderer *lipgloss.Renderer, options ...PrinterOption) Print {
	return NewPrinter(c, renderer, options...)
}

func (c Color) BufferPrinter(w io.Writer, options ...PrinterOption) Print {
	return c.Printer(Renderer(w), options...)
}

func (c Color) TextStyle(style lipgloss.Style) lipgloss.Style {
	if c.noColor {
		return style
	}

	return style.
		Foreground(c.Text)
}

func (c Color) TextEmphasisStyle(style lipgloss.Style) lipgloss.Style {
	if c.noColor {
		return style
	}

	return style.
		Foreground(c.TextEmphasis).
		Background(c.Background).
		Bold(true).
		BorderForeground(c.Border)
}

func (c Color) BoxStyles(header, body lipgloss.Style) Box {
	return Box{
		Header: header.
			Align(lipgloss.Center, lipgloss.Center).
			Border(headerBorder).
			BorderForeground(c.Border).
			PaddingBottom(1).
			PaddingTop(1).
			Inherit(c.TextEmphasisStyle(header)),

		Body: body.
			Align(lipgloss.Left).
			Border(bodyBorder).
			BorderForeground(c.Border).
			BorderTop(false).
			Padding(1),
	}
}
