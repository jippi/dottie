package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jippi/dottie/pkg/tui"
	"github.com/urfave/cli/v3"
)

var printCommand = &cli.Command{
	Name:   "print",
	Usage:  "Print environment variables",
	Before: setup,
	Action: func(_ context.Context, _ *cli.Command) error {
		settings.Interpolate = true

		if settings.ShowColors {
			tui.Theme.Dark.Printer(tui.Renderer(os.Stdout)).Println(env.Render(*settings))
		} else {
			fmt.Println(env.Render(*settings))
		}

		return nil
	},
}
