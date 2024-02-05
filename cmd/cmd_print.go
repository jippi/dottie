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
		settings.Apply(render.WithInterpolation(true))

		fmt.Println(render.NewRenderer(*settings).Statement(env).String())

		return nil
	},
}
