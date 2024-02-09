package disable

import (
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "disable KEY",
	Short: "Disable (comment) a KEY if it exists",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing required argument: KEY")
		}
		key := args[0]

		env, _, err := shared.Setup(cmd.Flags())
		if err != nil {
			return err
		}

		existing := env.Get(key)
		if existing == nil {
			return fmt.Errorf("Could not find KEY [%s]", key)
		}

		existing.Disable()

		return pkg.Save(cmd.Flag("file").Value.String(), env)
	},
}
