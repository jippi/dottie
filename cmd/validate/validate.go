package validate

import (
	"errors"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Validate an .env file",
		GroupID: "output",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}

	cmd.Flags().StringSlice("exclude-prefix", []string{}, "Exclude KEY with this prefix")
	cmd.Flags().StringSlice("ignore-rule", []string{}, "Ignore this validation rule (e.g. 'dir')")

	shared.BoolWithInverse(cmd, "fix", true, "Guide the user to fix supported validation errors", "Do not guide the user to fix supported validation errors")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	document, err := pkg.Load(cmd.Context(), filename)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	//
	// Interpolate
	//

	if err = document.InterpolateAll(cmd.Context()); err != nil {
		return fmt.Errorf("failed to interpolate file: %w", err)
	}

	//
	// Validate
	//

	var (
		excludedPrefixes = shared.StringSliceFlag(cmd.Flags(), "exclude-prefix")
		ignoreRules      = shared.StringSliceFlag(cmd.Flags(), "ignore-rule")
		stderr           = tui.StderrFromContext(cmd.Context())
		selectors        = []ast.Selector{
			ast.ExcludeDisabledAssignments,
		}
	)

	for _, filter := range excludedPrefixes {
		selectors = append(selectors, ast.ExcludeKeyPrefix(filter))
	}

	validationErrors, errs := document.Validate(cmd.Context(), selectors, ignoreRules)
	if errs != nil {
		return errs
	}

	if len(validationErrors) == 0 {
		stderr.Success().Box("No validation errors found")

		return nil
	}

	attemptFixOfValidationError := shared.BoolWithInverseValue(cmd.Flags(), "fix")

	danger := stderr.Danger()
	danger.Box(fmt.Sprintf("%d validation errors found", len(validationErrors)))
	danger.Println()

	for _, errIsh := range validationErrors {
		stderr.NoColor().Println(
			validation.Explain(
				cmd.Context(),
				document,
				errIsh,
				errIsh.Assignment,
				attemptFixOfValidationError,
				true,
			))
	}

	//
	// Validate file again, in case some of the fixers from before fixed them
	//

	document, err = pkg.Load(cmd.Context(), filename)
	if err != nil {
		return fmt.Errorf("failed to reload .env file: %w", err)
	}

	newRes, errs := document.Validate(cmd.Context(), selectors, ignoreRules)
	if errs != nil {
		return errs
	}

	if len(newRes) == 0 {
		stderr.Success().Println("All validation errors fixed")

		return nil
	}

	diff := len(validationErrors) - len(newRes)
	if diff > 0 {
		stderr.Warning().
			Box(
				fmt.Sprintf("%d validation errors left", len(newRes)),
				stderr.Success().Sprintf("%d validation errors was fixed", diff),
			)
	}

	return errors.New("Validation failed")
}
