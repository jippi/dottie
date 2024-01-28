package main

import (
	"context"
	"fmt"

	"dotfedi/pkg"
	"dotfedi/pkg/ast"

	"github.com/urfave/cli/v3"
)

var setCommand = &cli.Command{
	Name:      "set",
	Usage:     "Set a key/value pair",
	Before:    setup,
	ArgsUsage: "key value",
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
		value := cmd.Args().Get(1)

		assignment := env.Get(key)
		if assignment == nil {
			if cmd.Bool("error-if-missing") {
				return fmt.Errorf("Key [%s] does not exists", key)
			}
		}

		assignment.Value = value
		assignment.Commented = cmd.Bool("commented")
		assignment.SetQuote(cmd.String("quote-style"))

		if comments := cmd.StringSlice("comment"); len(comments) > 0 {
			slice := make([]*ast.Comment, 0)
			for _, v := range comments {
				slice = append(slice, ast.NewComment(v))
			}
			assignment.Comments = slice
		}

		return pkg.Save(cmd.String("file"), env)
	},
}
