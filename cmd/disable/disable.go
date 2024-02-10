package disable

import (
	"errors"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:               "disable KEY",
	Short:             "Disable (comment out) a KEY if it exists",
	ValidArgsFunction: shared.NewCompleter().WithHandlers(render.ExcludeDisabledAssignments).Get(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Missing required argument: KEY")
		}

		key := args[0]

		filename := cmd.Flag("file").Value.String()

		env, err := pkg.Load(filename)
		if err != nil {
			return err
		}

		existing := env.Get(key)
		if existing == nil {
			return fmt.Errorf("Could not find KEY [%s]", key)
		}

		existing.Disable()

		return pkg.Save(filename, env)
	},
}
