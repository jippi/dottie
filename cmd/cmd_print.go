package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var printCommand = &cli.Command{
	Name:   "print",
	Usage:  "Print environment variables",
	Before: setup,
	Action: func(_ context.Context, _ *cli.Command) error {
		res := env.RenderWithFilter(settings)
		fmt.Println(string(res))

		return nil
	},
}
