package value

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
		Use:               "value KEY",
		Short:             "Print value of a env key if it exists",
		GroupID:           "output",
		ValidArgsFunction: shared.NewCompleter().WithHandlers(render.ExcludeDisabledAssignments).Get(),
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
				return fmt.Errorf("Key [%s] does not exists", key)
			}

			if !existing.Enabled && !shared.BoolFlag(cmd.Flags(), "include-commented") {
				return fmt.Errorf("Key [%s] exists, but is commented out - use [--include-commented] to include it", key)
			}

			warn, err := env.InterpolateStatement(existing)
			if warn != nil {
				tui.Theme.Warning.StderrPrinter().Printfln("%+v", warn)
			}
			if err != nil {
				return err
			}

			fmt.Println(existing.Interpolated)

			return nil
		},
	}
}
