package main

import (
	"context"

	"github.com/jippi/dottie/pkg"
	"github.com/urfave/cli/v3"
)

var formatCmd = &cli.Command{
	Name:   "fmt",
	Usage:  "Format the file",
	Before: setup,
	Action: func(_ context.Context, cmd *cli.Command) error {
		pkg.Save(cmd.String("file"), env)

		return nil
	},
}
