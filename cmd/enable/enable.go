package enable

import (
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "enable",
	Short: "Enable (uncomment) a KEY if it exists",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing required argument: KEY")
		}

		env, _, err := shared.Setup(cmd.Flags())
		if err != nil {
			return err
		}

		key := args[0]

		existing := env.Get(key)
		if existing == nil {
			return fmt.Errorf("Could not find KEY [%s]", key)
		}

		existing.Enable()

		return pkg.Save(cmd.Flag("file").Value.String(), env)
	},
}
