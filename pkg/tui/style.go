package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type styleIdentifier int

const (
	Danger styleIdentifier = 1 << iota
	Dark
	Info
	Light
	NoColor
	Primary
	Secondary
	Success
	Warning
)

type Style struct {
	textColor         lipgloss.AdaptiveColor
	textStyle         lipgloss.Style
	textEmphasisColor lipgloss.AdaptiveColor
	textEmphasisStyle lipgloss.Style
	backgroundColor   lipgloss.AdaptiveColor
	borderColor       lipgloss.AdaptiveColor
}

func NewStyle(baseColor lipgloss.Color) Style {
	base := ColorToHex(baseColor)

	style := Style{
		textColor: lipgloss.AdaptiveColor{
			Light: TransformColor(base, "", 0),
			Dark:  TransformColor(base, "tint", 0.4),
		},
		textEmphasisColor: lipgloss.AdaptiveColor{
			Light: TransformColor(base, "shade", 0.6),
			Dark:  TransformColor(base, "tint", 0.4),
		},
		backgroundColor: lipgloss.AdaptiveColor{
			Light: TransformColor(base, "tint", 0.8),
			Dark:  TransformColor(base, "shade", 0.8),
		},
		borderColor: lipgloss.AdaptiveColor{
			Light: TransformColor(base, "tint", 0.6),
			Dark:  TransformColor(base, "shade", 0.4),
		},
	}

	style.textStyle = lipgloss.
		NewStyle().
		Foreground(style.textColor)

	style.textEmphasisStyle = lipgloss.
		NewStyle().
		Bold(true).
		Foreground(style.textEmphasisColor).
		Background(style.backgroundColor).
		BorderForeground(style.borderColor)

	return style
}

func NewStyleWithoutColor() Style {
	// Since all lipgloss.Styles are non-pointers, they are by default an empty / unstyled version of themselves
	return Style{}
}

func (style Style) NewPrinter(renderer *lipgloss.Renderer, options ...PrinterOption) Printer {
	return NewPrinter(style, renderer, options...)
}

func (style Style) TextStyle() lipgloss.Style {
	return style.textStyle
}

func (style Style) TextEmphasisStyle() lipgloss.Style {
	return style.textEmphasisStyle
}

func (style Style) BoxHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Border(headerBorder).
		BorderForeground(style.borderColor).
		PaddingBottom(1).
		PaddingTop(1).
		Inherit(style.TextEmphasisStyle())
}

func (style Style) BoxBody() lipgloss.Style {
	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Border(bodyBorder).
		BorderForeground(style.borderColor).
		BorderTop(false).
		Padding(1)
}
