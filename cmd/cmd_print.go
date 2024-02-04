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

		if settings.ShowPretty {
			fmt.Println(render.NewFormattedPresenter(*settings).Document(env, *settings))

			return nil
		}

		fmt.Println(render.NewDirectPresenter(*settings).Document(env, *settings))

		return nil
	},
}
