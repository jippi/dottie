package enable

import (
	"errors"
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
		GroupID:           "manipulate",
		ValidArgsFunction: shared.NewCompleter().WithHandlers(render.ExcludeActiveAssignments).Get(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("Missing required argument: KEY")
			}

			filename := cmd.Flag("file").Value.String()

			env, err := pkg.Load(filename)
			if err != nil {
				return err
			}

			key := args[0]

			existing := env.Get(key)
			if existing == nil {
				return fmt.Errorf("Could not find KEY [%s]", key)
			}

			stdout, stderr := tui.WritersFromContext(cmd.Context())

			if existing.Enabled {
				stderr.Warning().Printfln("WARNING: The key [%s] is already enabled", key)
			}

			existing.Enable()

			if err := pkg.Save(cmd.Context(), filename, env); err != nil {
				return fmt.Errorf("could not save file: %w", err)
			}

			stdout.Success().Printfln("Key [%s] was successfully enabled", key)

			return nil
		},
	}
}
