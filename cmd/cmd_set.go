package main

import (
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"

	"github.com/urfave/cli/v3"
)

var setCommand = &cli.Command{
	Name:      "set",
	Usage:     "Set a key/value pair",
	Before:    setup,
	ArgsUsage: "KEY VALUE",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "commented",
			OnlyOnce: true,
		},
		&cli.BoolFlag{
			Name:     "error-if-missing",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "group",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "before",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "after",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "quote-style",
			Usage:    "single|double|none",
			OnlyOnce: true,
		},
		&cli.StringSliceFlag{
			Name: "comment",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		key := cmd.Args().Get(0)
		if len(key) == 0 {
			return fmt.Errorf("Missing required argument: KEY")
		}

		options := ast.SetOptions{
			ErrorIfMissing: cmd.Bool("error-if-missing"),
			Before:         cmd.String("before"),
			Group:          cmd.String("group"),
			Comments:       cmd.StringSlice("comment"),
		}

		assignment := &ast.Assignment{
			Name:    key,
			Literal: cmd.Args().Get(1),
			Active:  !cmd.Bool("commented"),
		}

		assignment.SetQuote(cmd.String("quote-style"))

		_, err := env.Set(assignment, options)
		if err != nil {
			return err
		}

		return pkg.Save(cmd.String("file"), env)
	},
}
