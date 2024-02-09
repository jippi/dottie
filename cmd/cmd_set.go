package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"

	"github.com/urfave/cli/v3"
)

var setCommand = &cli.Command{
	Name:      "set",
	Usage:     "Set/update one or multiple key=value pairs",
	UsageText: "set KEY=VALUE [KEY=VALUE ...]",
	Before:    setup,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "disabled",
			Usage:    "Set/change the flag to be disabled (commented out)",
			Value:    false,
			OnlyOnce: true,
		},
		&cli.BoolFlag{
			Name:     "validate",
			Usage:    "Validate the VALUE input before saving the file",
			Value:    true,
			OnlyOnce: true,
		},
		&cli.BoolFlag{
			Name:     "error-if-missing",
			Usage:    "Exit with an error if the KEY does not exists in the .env file already",
			Value:    false,
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "group",
			Usage:    "The (optional) group name to add the KEY=VALUE pair under",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "before",
			Usage:    "If the key doesn't exist, add it to the file *before* this KEY",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "after",
			Usage:    "If the key doesn't exist, add it to the file *after* this KEY",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "quote-style",
			Usage:    "[single | double | none]",
			Value:    "double",
			OnlyOnce: true,
			Validator: func(s string) error {
				if token.QuoteFromString(s) > 0 {
					return nil
				}

				return fmt.Errorf("must be one of [single | double | none]")
			},
		},
		&cli.StringSliceFlag{
			Name:  "comment",
			Usage: "Set one or multiple lines of comments to the KEY=VALUE pair",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 {
			return fmt.Errorf("Missing required argument: KEY=VALUE")
		}

		options := ast.UpsertOptions{
			InsertBefore:   cmd.String("before"),
			Comments:       cmd.StringSlice("comment"),
			ErrorIfMissing: cmd.Bool("error-if-missing"),
			Group:          cmd.String("group"),
			SkipValidation: !cmd.Bool("validate"),
		}

		for _, stringPair := range cmd.Args().Slice() {
			pairSlice := strings.SplitN(stringPair, "=", 2)
			if len(pairSlice) != 2 {
				return fmt.Errorf("expected KEY=VALUE pair, missing '='")
			}

			key := pairSlice[0]
			value := pairSlice[1]

			assignment := &ast.Assignment{
				Name:    key,
				Literal: value,
				// by default we take the user input and assume its interpolated,
				// it will be interpolated inside (*Document).Set if applicable
				Interpolated: value,
				Active:       !cmd.Bool("disabled"),
				Quote:        token.QuoteFromString(cmd.String("quote-style")),
			}

			//
			// Upsert key
			//

			assignment, err := env.Upsert(assignment, options)
			if err != nil {
				validation.Explain(env, validation.NewError(assignment, err))

				return fmt.Errorf("failed to upsert the key/value pair")
			}

			tui.Theme.Success.StderrPrinter().Println("Key was successfully upserted")
		}

		//
		// Save file
		//

		if err := pkg.Save(cmd.String("file"), env); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		tui.Theme.Success.StderrPrinter().Println("File was successfully saved")

		return nil
	},
}
