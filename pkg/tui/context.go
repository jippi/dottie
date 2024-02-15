package tui

import (
	"context"
	"io"
)

type printerContextKey int

const (
	Stdout printerContextKey = iota
	Stderr printerContextKey = iota
)

func CreateContext(ctx context.Context, stdout, stderr io.Writer) context.Context {
	ctx = context.WithValue(ctx, Stdout, Theme.Printer(stdout))
	ctx = context.WithValue(ctx, Stderr, Theme.Printer(stderr))

	return ctx
}

func FromContext(ctx context.Context, key printerContextKey) ThemePrinter {
	return ctx.Value(key).(ThemePrinter)
}

func ColorFromContext(ctx context.Context, key printerContextKey, color colorType) Printer {
	return ctx.Value(key).(ThemePrinter).Color(color)
}

func PrintersFromContext(ctx context.Context) (ThemePrinter, ThemePrinter) {
	return ctx.Value(Stdout).(ThemePrinter), ctx.Value(Stderr).(ThemePrinter)
}
