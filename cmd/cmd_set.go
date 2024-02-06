package main

import (
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"

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
		&cli.BoolFlag{
			Name:     "skip-validation",
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

		options := ast.UpsertOptions{
			InsertBefore:   cmd.String("before"),
			Comments:       cmd.StringSlice("comment"),
			ErrorIfMissing: cmd.Bool("error-if-missing"),
			Group:          cmd.String("group"),
			SkipValidation: cmd.Bool("skip-validation"),
		}

		assignment := &ast.Assignment{
			Name:    key,
			Literal: cmd.Args().Get(1),
			// by default we take the user input and assume its interpolated,
			// it will be interpolated inside (*Document).Set if applicable
			Interpolated: cmd.Args().Get(1),
			Active:       !cmd.Bool("commented"),
		}

		assignment.SetQuote(cmd.String("quote-style"))

		// Upsert key

		assignment, err := env.Upsert(assignment, options)
		if err != nil {
			validation.Explain(env, validation.NewError(assignment, err))

			return fmt.Errorf("failed to upsert the key/value pair")
		}

		tui.Theme.Success.StderrPrinter().Println("Key was successfully upserted")

		if err := pkg.Save(cmd.String("file"), env); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		tui.Theme.Success.StderrPrinter().Println("File was successfully saved")

		return nil
	},
}
