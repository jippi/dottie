package disable

import (
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "disable KEY",
		Short:             "Disable (comment out) a KEY if it exists",
		GroupID:           "manipulate",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: shared.NewCompleter().WithHandlers(render.ExcludeDisabledAssignments).Get(),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := cmd.Flags().Arg(0)

			filename := cmd.Flag("file").Value.String()

			env, err := pkg.Load(filename)
			if err != nil {
				return err
			}

			assignment := env.Get(key)
			if assignment == nil {
				return fmt.Errorf("Could not find KEY [ %s ]", key)
			}

			if !assignment.Enabled {
				tui.MaybePrintWarnings(cmd.Context(), fmt.Errorf("The key [ %s ] is already disabled", key))

				return nil
			}

			assignment.Disable()

			if err := pkg.Save(cmd.Context(), filename, env); err != nil {
				return fmt.Errorf("could not save file: %w", err)
			}

			tui.StdoutFromContext(cmd.Context()).
				Success().
				Printfln("Key [ %s ] was successfully disabled", key)

			return nil
		},
	}
}
