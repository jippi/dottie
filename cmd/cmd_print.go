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

		fmt.Println(env.Render(*settings))

		return nil
	},
}
