package main

import (
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg/render"
	"github.com/urfave/cli/v3"
)

var printCommand = &cli.Command{
	Name:   "print",
	Usage:  "Print environment variables",
	Before: setup,
	Action: func(_ context.Context, _ *cli.Command) error {
		settings.Interpolate = true

		var handlers []render.Handler

		if settings.ShowPretty {
			handlers = append(handlers, render.FormatHandler)
		}

		fmt.Println(render.NewRenderer(*settings, handlers...).Document(env))

		return nil
	},
}
