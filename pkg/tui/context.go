package tui

import (
	"context"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type fileDescriptorKey int

const (
	Stdout fileDescriptorKey = iota
	Stderr
)

type contextKey int

const (
	themeContextValue contextKey = iota
	colorProfileContextValue
)

func NewContext(ctx context.Context, stdout, stderr io.Writer) context.Context {
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

func ColorProfileFromContext(ctx context.Context) termenv.Profile {
	return ctx.Value(colorProfileContextValue).(termenv.Profile) //nolint:forcetypeassert
}

func WriterFromContext(ctx context.Context, descriptor fileDescriptorKey) Writer {
	return ctx.Value(descriptor).(Writer) //nolint:forcetypeassert
}

func WritersFromContext(ctx context.Context) (Writer, Writer) {
	return ctx.Value(Stdout).(Writer), ctx.Value(Stderr).(Writer) //nolint:forcetypeassert
}
