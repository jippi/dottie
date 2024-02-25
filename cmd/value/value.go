package value

import (
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
	slogctx "github.com/veqryn/slog-context"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "value KEY",
		Short:             "Print value of a env key if it exists",
		GroupID:           "output",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.NewCompleter().WithSelectors(ast.ExcludeDisabledAssignments).Get(),
		RunE:              runE,
	}

	cmd.Flags().Bool("literal", false, "Show literal value instead of interpolated")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	document, err := pkg.Load(cmd.Context(), filename)
	if err != nil {
		return err
	}

	key := cmd.Flags().Arg(0)

	assignment := document.Get(key)
	if assignment == nil {
		return fmt.Errorf("Key [ %s ] does not exists", key)
	}

	if !assignment.Enabled && !shared.BoolFlag(cmd.Flags(), "include-commented") {
		return fmt.Errorf("Key [ %s ] exists, but is commented out - use [--include-commented] to include it", key)
	}

	if ok, _ := cmd.Flags().GetBool("literal"); ok {
		out, err := assignment.Unquote(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Fprint(cmd.OutOrStdout(), out)

		return nil
	}

	if err := document.InterpolateStatement(cmd.Context(), assignment); err != nil {
		return err
	}

	slogctx.Info(cmd.Context(), "value.assignment.Interpolated", tui.StringDump("assignment.Interpolated", assignment.Interpolated))

	fmt.Fprint(cmd.OutOrStdout(), assignment.Interpolated)

	return nil
}
