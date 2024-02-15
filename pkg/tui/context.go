package tui

import (
	"context"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type printerContextKey int

const (
	Stdout printerContextKey = iota
	Stderr
)

type themeContextKey int

const (
	themeContextValue themeContextKey = iota
	colorProfileContextValue
)

func CreateContext(ctx context.Context, stdout, stderr io.Writer) context.Context {
	theme := NewTheme()

	stdoutOutput := lipgloss.NewRenderer(stdout, termenv.WithColorCache(true))
	stderrOutput := lipgloss.NewRenderer(stderr, termenv.WithColorCache(true))

	ctx = context.WithValue(ctx, themeContextValue, theme)
	ctx = context.WithValue(ctx, colorProfileContextValue, stdoutOutput.ColorProfile())
	ctx = context.WithValue(ctx, Stdout, theme.NewWriter(stdoutOutput))
	ctx = context.WithValue(ctx, Stderr, theme.NewWriter(stderrOutput))

	return ctx
}

func ThemeFromContext(ctx context.Context) Theme {
	return ctx.Value(themeContextValue).(Theme) //nolint:forcetypeassert
}

func ColorProfile(ctx context.Context) termenv.Profile {
	return ctx.Value(colorProfileContextValue).(termenv.Profile) //nolint:forcetypeassert
}

func WriterFromContext(ctx context.Context, key printerContextKey) Writer {
	return ctx.Value(key).(Writer) //nolint:forcetypeassert
}

func ColorPrinterFromContext(ctx context.Context, key printerContextKey, color colorType) Printer {
	return ctx.Value(key).(*Writer).Color(color) //nolint:forcetypeassert
}

func PrintersFromContext(ctx context.Context) (Writer, Writer) {
	return ctx.Value(Stdout).(Writer), ctx.Value(Stderr).(Writer) //nolint:forcetypeassert
}
