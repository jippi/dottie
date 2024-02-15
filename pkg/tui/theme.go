package tui

import (
	"context"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type Theme struct {
	Danger    Style
	Dark      Style
	Info      Style
	Light     Style
	NoColor   Style
	Primary   Style
	Secondary Style
	Success   Style
	Warning   Style
}

func NewTheme() Theme {
	theme := Theme{}

	theme.Danger = NewStyle(Red)
	theme.Info = NewStyle(Cyan)
	theme.Light = NewStyle(Gray300)
	theme.NoColor = NewStyleWithoutColor()
	theme.Primary = NewStyle(Blue)
	theme.Secondary = NewStyle(Gray600)
	theme.Success = NewStyle(Green)
	theme.Warning = NewStyle(Yellow)

	theme.Dark = NewStyle(Gray700)
	theme.Dark.textEmphasisColor.Dark = ColorToHex(Gray300)
	theme.Dark.backgroundColor.Dark = "#1a1d20"
	theme.Dark.borderColor.Dark = ColorToHex(Gray800)

	return theme
}

func (theme Theme) NewWriter(writer *lipgloss.Renderer) Writer {
	return Writer{
		writer: writer,
		theme:  theme,
		cache:  make(map[style]Printer),
	}
}

func NewWriter(ctx context.Context, writer io.Writer) Writer {
	var options []termenv.OutputOption

	// If the primary color profile is in color mode, enforce TTY to keep coloring on
	if profile := ColorProfileFromContext(ctx); profile != termenv.Ascii {
		options = append(
			options,
			termenv.WithTTY(true),
			termenv.WithProfile(profile),
		)
	}

	return ThemeFromContext(ctx).NewWriter(lipgloss.NewRenderer(writer, options...))
}
