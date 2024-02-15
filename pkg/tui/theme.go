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

	theme.Danger = NewStyleFromColor(Red)
	theme.Info = NewStyleFromColor(Cyan)
	theme.Light = NewStyleFromColor(Gray300)
	theme.Primary = NewStyleFromColor(Blue)
	theme.Secondary = NewStyleFromColor(Gray600)
	theme.Success = NewStyleFromColor(Green)
	theme.Warning = NewStyleFromColor(Yellow)
	theme.NoColor = NewStyleWithoutColor()

	dark := NewStyleConfig(Gray700)
	dark.TextEmphasis.Dark = ColorToHex(Gray300)
	dark.Background.Dark = "#1a1d20"
	dark.Border.Dark = ColorToHex(Gray800)

	theme.Dark = NewStyle(dark)

	return theme
}
