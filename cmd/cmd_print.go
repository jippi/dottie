package main

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v3"
)

var printCommand = &cli.Command{
	Name:   "print",
	Usage:  "Print environment variables",
	Before: setup,
	Action: func(_ context.Context, _ *cli.Command) error {
		spew.Config.DisablePointerMethods = true
		spew.Config.DisableMethods = true

		// spew.Dump(env)
		// return nil

		res := env.RenderWithFilter(*settings)
		fmt.Println(string(res))

		return nil
	},
}
