package tui

import (
	"context"
	"io"
	"log/slog"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/reugn/pkgslog"
	slogmulti "github.com/samber/slog-multi"
	slogctx "github.com/veqryn/slog-context"
	slogdedup "github.com/veqryn/slog-dedup"
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
	ctx = NewContextWithoutLogger(ctx, stdout, stderr)
	ctx = slogctx.NewCtx(
		ctx,
		slog.New(
			slogmulti.
				Pipe(
					func(next slog.Handler) slog.Handler {
						return pkgslog.NewPackageHandler(next, packageLogLevels())
					},
				).
				Pipe(
					slogctx.NewMiddleware(&slogctx.HandlerOptions{}),
				).
				Pipe(
					slogdedup.NewOverwriteMiddleware(&slogdedup.OverwriteHandlerOptions{
						ResolveKey: slogdedup.KeepIfBuiltinKeyConflict,
					}),
				).
				Handler(
					logHandler(stderr),
				),
		),
	)

	return ctx
}

func NewContextWithoutLogger(ctx context.Context, stdout, stderr io.Writer) context.Context {
	theme := NewTheme()

	stdoutOutput := lipgloss.NewRenderer(stdout, termenv.WithColorCache(true))
	stderrOutput := lipgloss.NewRenderer(stderr, termenv.WithColorCache(true))

	ctx = context.WithValue(ctx, themeContextValue, theme)
	ctx = context.WithValue(ctx, colorProfileContextValue, stdoutOutput.ColorProfile())
	ctx = context.WithValue(ctx, Stdout, theme.Writer(stdoutOutput))
	ctx = context.WithValue(ctx, Stderr, theme.Writer(stderrOutput))

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

func StdoutFromContext(ctx context.Context) Writer {
	return WriterFromContext(ctx, Stdout)
}

func StderrFromContext(ctx context.Context) Writer {
	return WriterFromContext(ctx, Stderr)
}

func WritersFromContext(ctx context.Context) (Writer, Writer) {
	return StdoutFromContext(ctx), StderrFromContext(ctx)
}
