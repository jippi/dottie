package main

import (
	"context"

	"dotfedi/pkg"

	"github.com/urfave/cli/v3"
)

var disableCommand = &cli.Command{
	Name:      "disable",
	Usage:     "Comment/disable a key if it exists",
	Before:    setup,
	ArgsUsage: "KEY",
	Action: func(_ context.Context, cmd *cli.Command) error {
		key := cmd.Args().Get(0)
		existing := env.Get(key)
		existing.Active = true

		return pkg.Save(cmd.String("file"), env)
	},
}
