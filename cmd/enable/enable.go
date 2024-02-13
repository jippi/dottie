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

var Command = &cobra.Command{
	Use:               "enable KEY",
	Short:             "Enable (uncomment) a KEY if it exists",
	GroupID:           "manipulate",
	ValidArgsFunction: shared.NewCompleter().WithHandlers(render.ExcludeActiveAssignments).Get(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Missing required argument: KEY")
		}

		filename := cmd.Flag("file").Value.String()

		env, warn, err := pkg.Load(filename)
		if warn != nil {
			tui.Theme.Warning.StderrPrinter().Println(warn)
		}
		if err != nil {
			return err
		}

		key := args[0]

		existing := env.Get(key)
		if existing == nil {
			return fmt.Errorf("Could not find KEY [%s]", key)
		}

		existing.Enable()

		return pkg.Save(filename, env)
	},
}
