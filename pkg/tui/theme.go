package tui

import (
	"context"
	"io"

	lipgloss "charm.land/lipgloss/v2"
)

type Theme struct {
	styles map[styleIdentifier]Style
}

func NewTheme() Theme {
	theme := Theme{}
	theme.styles = make(map[styleIdentifier]Style)
	theme.styles[Danger] = NewStyle(Red)
	theme.styles[Info] = NewStyle(Cyan)
	theme.styles[Light] = NewStyle(Gray300)
	theme.styles[NoColor] = NewStyleWithoutColor()
	theme.styles[Primary] = NewStyle(Blue)
	theme.styles[Secondary] = NewStyle(Gray600)
	theme.styles[Success] = NewStyle(Green)
	theme.styles[Warning] = NewStyle(Yellow)

	dark := NewStyle(Gray700)
	dark.textEmphasisColor.Dark = lipgloss.Color(Gray300)
	dark.backgroundColor.Dark = lipgloss.Color("#1a1d20")
	dark.borderColor.Dark = lipgloss.Color(Gray800)

	theme.styles[Dark] = dark

	return theme
}

func (theme Theme) Style(id styleIdentifier) Style {
	return theme.styles[id]
}

func (theme Theme) Writer(w io.Writer) Writer {
	return Writer{
		writer: w,
		theme:  theme,
		cache:  make(map[styleIdentifier]Printer),
	}
}

func NewWriter(ctx context.Context, writer io.Writer) Writer {
	return ThemeFromContext(ctx).Writer(writer)
}
