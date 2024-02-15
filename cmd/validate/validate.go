package validate

import (
	"errors"
	"fmt"
	"os"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Validate an .env file",
		GroupID: "output",
		RunE:    runE,
	}

	cmd.Flags().StringSlice("exclude-prefix", []string{}, "Exclude KEY with this prefix")
	cmd.Flags().StringSlice("ignore-rule", []string{}, "Ignore this validation rule (e.g. 'dir')")

	shared.BoolWithInverse(cmd, "fix", true, "Guide the user to fix supported validation errors", "Do not guide the user to fix supported validation errors")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	env, err := pkg.Load(filename)
	if err != nil {
		return err
	}

	//
	// Build filters
	//

	fix := shared.BoolWithInverseValue(cmd.Flags(), "fix")
	ignoreRules, _ := cmd.Flags().GetStringSlice("ignore-rule")

	handlers := []render.Handler{}
	handlers = append(handlers, render.ExcludeDisabledAssignments)

	slice, _ := cmd.Flags().GetStringSlice("exclude-prefix")
	for _, filter := range slice {
		handlers = append(handlers, render.ExcludeKeyPrefix(filter))
	}

	//
	// Interpolate
	//

	warn, err := env.InterpolateAll()

	if warn != nil {
		tui.Theme.Warning.StderrPrinter().Printfln("%+v", warn)
	}

	if err != nil {
		return err
	}

	//
	// Validate
	//

	res := validation.Validate(env, handlers, ignoreRules)
	if len(res) == 0 {
		tui.Theme.Success.StderrPrinter().Box("No validation errors found")

		return nil
	}

	stderr := tui.Theme.Danger.StderrPrinter()
	stderr.Box(fmt.Sprintf("%d validation errors found", len(res)))
	stderr.Println()

	for _, errIsh := range res {
		fmt.Fprintln(os.Stderr, validation.Explain(env, errIsh, errIsh, fix, true))
	}

	//
	// Validate file again, in case some of the fixers from before fixed them
	//

	env, err = pkg.Load(cmd.Flag("file").Value.String())
	if err != nil {
		return fmt.Errorf("failed to reload .env file: %w", err)
	}

	newRes := validation.Validate(env, handlers, ignoreRules)
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
}
