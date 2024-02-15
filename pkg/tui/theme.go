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

func (theme Theme) Printer(writer *lipgloss.Renderer) ThemeWriter {
	return ThemeWriter{
		writer: writer,
		theme:  theme,
		cache:  make(map[colorType]Printer),
	}
}

func (theme Theme) WriterPrinter(ctx context.Context, writer io.Writer) ThemeWriter {
	options := []termenv.OutputOption{}

	if ColorProfile(ctx) != termenv.Ascii {
		options = append(options, termenv.WithTTY(true))
	}

	return theme.Printer(lipgloss.NewRenderer(writer, options...))
}

func NewTheme() Theme {
	theme := Theme{}

	theme.Danger = NewStyle(Red)
	theme.Info = NewStyle(Cyan)
	theme.Light = NewStyle(Gray300)
	theme.Primary = NewStyle(Blue)
	theme.Secondary = NewStyle(Gray600)
	theme.Success = NewStyle(Green)
	theme.Warning = NewStyle(Yellow)
	theme.NoColor = NewStyleWithoutColor()

	theme.Dark = NewStyle(Gray700)
	theme.Dark.TextEmphasis.Dark = ColorToHex(Gray300)
	theme.Dark.Background.Dark = "#1a1d20"
	theme.Dark.Border.Dark = ColorToHex(Gray800)

	return theme
}
