package set

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/ast/upsert"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
	slogctx "github.com/veqryn/slog-context"
	"go.uber.org/multierr"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set KEY=VALUE [KEY=VALUE ...]",
		Short:   "Set/update one or multiple key=value pairs",
		GroupID: "manipulate",
		Args:    cobra.MinimumNArgs(1),
		ValidArgsFunction: shared.NewCompleter().
			WithSuffixIsLiteral(true).
			WithHandlers(ast.ExcludeDisabledAssignments).
			Get(),
		RunE: runE,
	}

	shared.BoolWithInverse(cmd, "validate", true, "Validate the VALUE input before saving the file", "Do not validate the VALUE input before saving the file")

	cmd.Flags().Bool("disabled", false, "Set/change the flag to be disabled (commented out)")
	cmd.Flags().Bool("error-if-missing", false, "Exit with an error if the KEY does not exists in the .env file already")
	cmd.Flags().Bool("skip-if-exists", false, "If the already KEY exists, do not set or change any settings")
	cmd.Flags().Bool("skip-if-same", false, "If the already KEY exists, and it the value is identical, do not set or change any settings")

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

	document, err := pkg.Load(cmd.Context(), filename)
	if err != nil {
		return err
	}

	//
	// Initialize Upserter
	//

	upserter, err := upsert.New(
		document,
		upsert.WithGroup(shared.StringFlag(cmd.Flags(), "group")),
		upsert.EnableSettingIf(upsert.ErrorIfMissing, shared.BoolFlag(cmd.Flags(), "error-if-missing")),
		upsert.EnableSettingIf(upsert.SkipIfExists, shared.BoolFlag(cmd.Flags(), "skip-if-exists")),
		upsert.EnableSettingIf(upsert.SkipIfSame, shared.BoolFlag(cmd.Flags(), "skip-if-same")),
		upsert.EnableSettingIf(upsert.UpdateComments, cmd.Flag("comment").Changed),
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

	if err := upserter.ApplyOptions(upsert.EnableSettingIf(upsert.Validate, shared.BoolWithInverseValue(cmd.Flags(), "validate"))); err != nil {
		return fmt.Errorf("error configuring [--validate] flag: %w", err)
	}

	//
	// Loop arguments and place them
	//

	var (
		allErrors      error
		stdout, stderr = tui.WritersFromContext(cmd.Context())
	)

	var (
		argumentCounter  int
		skipNextArgument bool
	)

	for _, arg := range args {
		slogctx.Debug(cmd.Context(), "arg", tui.StringDump(arg))
	}

	for _, stringPair := range args {
		if skipNextArgument {
			skipNextArgument = false

			continue
		}

		var key, value string

		// KEY1=VALUE KEY2=VALUE [...]
		pairSlice := strings.SplitN(stringPair, "=", 2)
		if len(pairSlice) == 2 {
			argumentCounter++

			key = pairSlice[0]
			value = pairSlice[1]
		}

		// KEY1 VALUE1 KEY2 VALUE2
		if len(key) == 0 {
			key = args[argumentCounter]
			argumentCounter++

			if argumentCounter >= len(args) {
				allErrors = multierr.Append(allErrors, fmt.Errorf("Key [ %s ] Error: expected [KEY VALUE] arguments pair, missing [ VALUE ] argument", stringPair))

				break
			}

			value = args[argumentCounter]
			argumentCounter++

			skipNextArgument = true
		}

		// Fail
		if len(key) == 0 {
			allErrors = multierr.Append(allErrors, fmt.Errorf("Key [ %s ] Error: expected KEY=VALUE pair, missing '='", stringPair))

			continue
		}

		assignment := &ast.Assignment{
			Name:         key,
			Enabled:      !shared.BoolFlag(cmd.Flags(), "disabled"),
			Literal:      value,
			Interpolated: value,
			Quote:        token.QuoteFromString(shared.StringFlag(cmd.Flags(), "quote-style")),
			Comments:     ast.NewCommentsFromSlice(shared.StringSliceFlag(cmd.Flags(), "comment")),
		}

		assignment.SetLiteral(cmd.Context(), value)

		//
		// Upsert the assignment
		//

		var skippedStatementWarning upsert.SkippedStatementError

		assignment, err := upserter.Upsert(cmd.Context(), assignment)

		switch {
		case errors.As(err, &skippedStatementWarning):
			stderr.Warning().Print("  ", key)
			stderr.Dark().Print(" was skipped: ")
			stderr.Dark().Println(skippedStatementWarning.Reason)

		case err != nil:
			stderr.NoColor().Println(validation.Explain(cmd.Context(), document, err, assignment, false, true))

			if shared.BoolWithInverseValue(cmd.Flags(), "validate") {
				allErrors = multierr.Append(allErrors, err)

				continue
			}
		}

		stdout.Success().Printfln("Key [ %s ] was successfully upserted", key)
	}

	if allErrors != nil {
		return fmt.Errorf("validation error: %+w", allErrors)
	}

	//
	// Save file
	//

	if err := pkg.Save(cmd.Context(), shared.StringFlag(cmd.Flags(), "file"), document); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	stdout.Success().Println("File was successfully saved")

	return nil
}
