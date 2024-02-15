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

	theme.Danger = NewColor(NewColorComponentConfig(Red))
	theme.Info = NewColor(NewColorComponentConfig(Cyan))
	theme.Light = NewColor(NewColorComponentConfig(Gray300))
	theme.Primary = NewColor(NewColorComponentConfig(Blue))
	theme.Secondary = NewColor(NewColorComponentConfig(Gray600))
	theme.Success = NewColor(NewColorComponentConfig(Green))
	theme.Warning = NewColor(NewColorComponentConfig(Yellow))
	theme.NoColor = NewStyle()

	dark := NewColorComponentConfig(Gray700)

	dark.TextEmphasis.Dark = ComponentColorConfig{
		Color: ColorToHex(Gray300),
	}

	dark.Background.Dark = ComponentColorConfig{
		Color: "#1a1d20",
	}

	dark.Border.Dark = ComponentColorConfig{
		Color: ColorToHex(Gray800),
	}

	theme.Dark = NewColor(dark)

	return theme
}
