package tui

import (
	"context"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
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
	dark.textEmphasisColor.Dark = ColorToHex(Gray300)
	dark.backgroundColor.Dark = "#1a1d20"
	dark.borderColor.Dark = ColorToHex(Gray800)

	theme.styles[Dark] = dark

	return theme
}

func (theme Theme) Style(id styleIdentifier) Style {
	return theme.styles[id]
}

func (theme Theme) Writer(renderer *lipgloss.Renderer) Writer {
	return Writer{
		renderer: renderer,
		theme:    theme,
		cache:    make(map[styleIdentifier]Printer),
	}
}

func NewWriter(ctx context.Context, writer io.Writer) Writer {
	var options []termenv.OutputOption

	// If the primary (stdout) color profile is in color mode (aka not ASCII),
	// force  TTY and color profile for the new renderer and writer
	if profile := ColorProfileFromContext(ctx); profile != termenv.Ascii {
		options = append(
			options,
			termenv.WithTTY(true),
			termenv.WithProfile(profile),
		)
	}

	return ThemeFromContext(ctx).Writer(lipgloss.NewRenderer(writer, options...))
}
