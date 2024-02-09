package validate

import (
	"errors"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "validate",
	Short: "Validate .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, _, err := shared.Setup(cmd.Flags())
		if err != nil {
			return err
		}

		res := validation.Validate(env)
		if len(res) == 0 {
			tui.Theme.Success.StderrPrinter().Box("No validation errors found")

			return nil
		}

		stderr := tui.Theme.Danger.StderrPrinter()
		stderr.Box(fmt.Sprintf("%d validation errors found", len(res)))
		stderr.Println()

		for _, errIsh := range res {
			validation.Explain(env, errIsh)
		}

		env, err = pkg.Load(cmd.Flag("file").Value.String())
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
