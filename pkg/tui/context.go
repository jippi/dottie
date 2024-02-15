package tui

import (
	"context"
	"io"
)

type printerContextKey int

const (
	Stdout printerContextKey = iota
	Stderr
)

type themeContextKey int

const (
	Theme themeContextKey = iota
)

func CreateContext(ctx context.Context, stdout, stderr io.Writer) context.Context {
	theme := NewTheme()

	ctx = context.WithValue(ctx, Theme, theme)
	ctx = context.WithValue(ctx, Stdout, theme.Printer(stdout))
	ctx = context.WithValue(ctx, Stderr, theme.Printer(stderr))

	return ctx
}

func ThemeFromContext(ctx context.Context) ThemeConfig {
	return ctx.Value(Theme).(ThemeConfig) //nolint:forcetypeassert
}

func PrinterFromContext(ctx context.Context, key printerContextKey) ThemePrinter {
	return ctx.Value(key).(ThemePrinter) //nolint:forcetypeassert
}

func ColorPrinterFromContext(ctx context.Context, key printerContextKey, color colorType) Printer {
	return ctx.Value(key).(ThemePrinter).Color(color) //nolint:forcetypeassert
}

func PrintersFromContext(ctx context.Context) (ThemePrinter, ThemePrinter) {
	return ctx.Value(Stdout).(ThemePrinter), ctx.Value(Stderr).(ThemePrinter) //nolint:forcetypeassert
}
