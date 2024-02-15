package set

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/ast/upsert"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set KEY=VALUE [KEY=VALUE ...]",
		Short:   "Set/update one or multiple key=value pairs",
		GroupID: "manipulate",
		ValidArgsFunction: shared.NewCompleter().
			WithSuffixIsLiteral(true).
			WithHandlers(render.ExcludeDisabledAssignments).
			Get(),
		RunE: runE,
	}

	shared.BoolWithInverse(cmd, "validate", true, "Validate the VALUE input before saving the file", "Do not validate the VALUE input before saving the file")

	cmd.Flags().Bool("disabled", false, "Set/change the flag to be disabled (commented out)")
	cmd.Flags().Bool("error-if-missing", false, "Exit with an error if the KEY does not exists in the .env file already")
	cmd.Flags().String("group", "", "The (optional) group name to add the KEY=VALUE pair under")
	cmd.Flags().String("before", "", "If the key doesn't exist, add it to the file *before* this KEY")
	cmd.Flags().String("after", "", "If the key doesn't exist, add it to the file *after* this KEY")
	cmd.Flags().String("quote-style", "double", "The quote style to use (single, double, none)")
	cmd.Flags().StringSlice("comment", nil, "Set one or multiple lines of comments to the KEY=VALUE pair")

	cmd.MarkFlagsMutuallyExclusive("before", "after", "group")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	document, err := pkg.Load(filename)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("Missing required argument: KEY=VALUE")
	}

	//
	// Initialize Upserter
	//

	upserter, err := upsert.New(
		document,
		upsert.WithGroup(shared.StringFlag(cmd.Flags(), "group")),
		upsert.WithSettingIf(upsert.ErrorIfMissing, shared.BoolFlag(cmd.Flags(), "error-if-missing")),
		upsert.WithSettingIf(upsert.UpdateComments, cmd.Flag("comment").Changed),
	)
	if err != nil {
		return fmt.Errorf("error setting up upserter: %w", err)
	}

	if err := upserter.ApplyOptions(upsert.WithPlacementInGroupIgnoringEmpty(upsert.AddBeforeKey, shared.StringFlag(cmd.Flags(), "before"))); err != nil {
		return fmt.Errorf("error in processing [--before] flag: %w", err)
	}

	if err := upserter.ApplyOptions(upsert.WithPlacementInGroupIgnoringEmpty(upsert.AddAfterKey, shared.StringFlag(cmd.Flags(), "after"))); err != nil {
		return fmt.Errorf("error in processing [--after] flag: %w", err)
	}

	if err := upserter.ApplyOptions(upsert.WithSettingIf(upsert.Validate, shared.BoolWithInverseValue(cmd.Flags(), "validate"))); err != nil {
		return fmt.Errorf("error configuring [--validate] flag: %w", err)
	}

	//
	// Loop arguments and place them
	//

	var (
		allErrors      error
		stdout, stderr = tui.PrintersFromContext(cmd.Context())
	)

	for _, stringPair := range args {
		pairSlice := strings.SplitN(stringPair, "=", 2)
		if len(pairSlice) != 2 {
			allErrors = multierr.Append(allErrors, fmt.Errorf("Key: '%s' Error: expected KEY=VALUE pair, missing '='", stringPair))

			continue
		}

		key := pairSlice[0]
		value := pairSlice[1]

		assignment := &ast.Assignment{
			Name:         key,
			Literal:      value,
			Interpolated: value,
			Enabled:      !shared.BoolFlag(cmd.Flags(), "disabled"),
			Quote:        token.QuoteFromString(shared.StringFlag(cmd.Flags(), "quote-style")),
			Comments:     ast.NewCommentsFromSlice(shared.StringSliceFlag(cmd.Flags(), "comment")),
		}

		//
		// Upsert the assignment
		//

		assignment, warnings, err := upserter.Upsert(assignment)
		if warnings != nil {
			stderr.Color(tui.Warning).Println("WARNING:", warnings)
		}

		if err != nil {
			z := validation.NewError(assignment, err)
			stderr.Color(tui.Neutral).Println(validation.Explain(document, z, z, false, true))

			if shared.BoolWithInverseValue(cmd.Flags(), "validate") {
				allErrors = multierr.Append(allErrors, err)

				continue
			}
		}

		stdout.Color(tui.Success).Printfln("Key [%s] was successfully upserted", key)
	}

	if allErrors != nil {
		return errors.New("validation error")
	}

	//
	// Save file
	//

	if err := pkg.Save(shared.StringFlag(cmd.Flags(), "file"), document); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	stdout.Color(tui.Success).Println("File was successfully saved")

	return nil
}
