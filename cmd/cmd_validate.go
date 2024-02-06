package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/urfave/cli/v3"
)

var validateCommand = &cli.Command{
	Name:   "validate",
	Usage:  "Validate .env file",
	Before: setup,
	Action: func(_ context.Context, cmd *cli.Command) error {
		res := validation.Validate(env)
		if len(res) == 0 {
			tui.Theme.Success.StderrPrinter().Box("No validation errors found")

			return nil
		}

		stderr := tui.Theme.Danger.StderrPrinter()
		stderr.Box(fmt.Sprintf("%d validation errors found", len(res)))
		stderr.Println()

		env, err := pkg.Load(cmd.String("file"))
		if err != nil {
			return fmt.Errorf("failed to reload .env file: %w", err)
		}

		newRes := validation.Validate(env)
		if len(newRes) == 0 {
			tui.Theme.Success.StderrPrinter().Println("All validation errors fixed")

			return nil
		}

		diff := len(res) - len(newRes)
		if diff > 0 {
			tui.Theme.Warning.StderrPrinter().Box(
				fmt.Sprintf("%d validation errors left", len(newRes)),
				tui.Theme.Success.StderrPrinter().Sprintf("%d validation errors was fixed", diff),
			)
		}

		return errors.New("Validation failed")
	},
}
