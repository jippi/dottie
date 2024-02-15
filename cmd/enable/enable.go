package enable

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
		Use:               "enable KEY",
		Short:             "Enable (uncomment) a KEY if it exists",
		Args:              cobra.ExactArgs(1),
		GroupID:           "manipulate",
		ValidArgsFunction: shared.NewCompleter().WithHandlers(render.ExcludeActiveAssignments).Get(),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := cmd.Flag("file").Value.String()

			env, err := pkg.Load(filename)
			if err != nil {
				return err
			}

			key := cmd.Flags().Arg(0)

			assignment := env.Get(key)
			if assignment == nil {
				return fmt.Errorf("Could not find KEY [%s]", key)
			}

			if assignment.Enabled {
				tui.StderrFromContext(cmd.Context()).
					Warning().
					Printfln("WARNING: The key [%s] is already enabled", key)
			}

			assignment.Enable()

			if err := pkg.Save(cmd.Context(), filename, env); err != nil {
				return fmt.Errorf("could not save file: %w", err)
			}

			tui.StdoutFromContext(cmd.Context()).
				Success().
				Printfln("Key [%s] was successfully enabled", key)

			return nil
		},
	}
}
