package fmt

import (
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "fmt",
	Short: "Format a .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, _, err := shared.Setup(cmd.Flags())
		if err != nil {
			return err
		}

		return pkg.Save(cmd.Flag("file").Value.String(), env)
	},
}
