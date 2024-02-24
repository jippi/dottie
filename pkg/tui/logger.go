package tui

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/golang-cz/devslog"
	"github.com/lmittmann/tint"
)

func ParseLogLevel(name string, fallback slog.Level) slog.Leveler {
	switch strings.ToUpper(name) {
	case "DEBUG":
		return slog.LevelDebug

	case "INFO":
		return slog.LevelInfo

	case "WARN":
		return slog.LevelWarn

	case "ERROR":
		return slog.LevelError

	default:
		return fallback
	}
}

func logHandler(out io.Writer) slog.Handler {
	logLevel := ParseLogLevel(os.Getenv("LOG_LEVEL"), slog.LevelInfo)

	if _, ok := os.LookupEnv("CI"); ok {
		return tint.NewHandler(
			out,
			&tint.Options{
				Level:     logLevel,
				AddSource: logLevel == slog.LevelDebug,
			},
		)
	}

	return devslog.NewHandler(
		out,
		&devslog.Options{
			SortKeys: true,
			HandlerOptions: &slog.HandlerOptions{
				Level:     logLevel,
				AddSource: logLevel == slog.LevelDebug,
			},
		},
	)
}

func StringDump(value string) slog.Attr {
	return slog.Group(
		"string_dump",
		slog.String("normal", value),
		QuotedString("quoted", value),
		UnicodeString("unicode", value),
	)
}

func QuotedString(key, value string) slog.Attr {
	return slog.String(key, fmt.Sprintf("%q", value))
}

func UnicodeString(key, value string) slog.Attr {
	return slog.String(key, fmt.Sprintf("%U", []rune(value)))
}
