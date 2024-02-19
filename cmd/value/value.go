package value

import (
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "value KEY",
		Short:             "Print value of a env key if it exists",
		GroupID:           "output",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.NewCompleter().WithHandlers(ast.ExcludeDisabledAssignments).Get(),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := cmd.Flag("file").Value.String()

			document, err := pkg.Load(filename)
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

			warnings, err := document.InterpolateStatement(assignment)
			tui.MaybePrintWarnings(cmd.Context(), warnings)
			if err != nil {
				return err
			}

			tui.StdoutFromContext(cmd.Context()).
				NoColor().
				Println(assignment.Interpolated)

			return nil
		},
	}
}
