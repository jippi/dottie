package tui

import (
	"io"

	lipgloss "charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
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
	textColor         compat.AdaptiveColor
	textStyle         lipgloss.Style
	textEmphasisColor compat.AdaptiveColor
	textEmphasisStyle lipgloss.Style
	backgroundColor   compat.AdaptiveColor
	borderColor       compat.AdaptiveColor
}

func NewStyle(baseColor string) Style {
	base := baseColor

	style := Style{
		textColor: compat.AdaptiveColor{
			Light: lipgloss.Color(TransformColor(base, "", 0)),
			Dark:  lipgloss.Color(TransformColor(base, "tint", 0.4)),
		},
		textEmphasisColor: compat.AdaptiveColor{
			Light: lipgloss.Color(TransformColor(base, "shade", 0.6)),
			Dark:  lipgloss.Color(TransformColor(base, "tint", 0.4)),
		},
		backgroundColor: compat.AdaptiveColor{
			Light: lipgloss.Color(TransformColor(base, "tint", 0.8)),
			Dark:  lipgloss.Color(TransformColor(base, "shade", 0.8)),
		},
		borderColor: compat.AdaptiveColor{
			Light: lipgloss.Color(TransformColor(base, "tint", 0.6)),
			Dark:  lipgloss.Color(TransformColor(base, "shade", 0.4)),
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

func (style Style) NewPrinter(w io.Writer, options ...PrinterOption) Printer {
	return NewPrinter(style, w, options...)
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
