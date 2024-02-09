package value

import (
	"fmt"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "value",
	Short: "Print value of a env key if it exists",
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
			return fmt.Errorf("Key [%s] does not exists", key)
		}

		if !existing.Active && !shared.BoolFlag(cmd.Flags(), "include-commented") {
			return fmt.Errorf("Key [%s] exists, but is commented out - use [--include-commented] to include it", key)
		}

		fmt.Println(existing.Interpolated)

		return nil
	},
}
