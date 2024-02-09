package groups

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "groups",
	Short: "Print groups found in the .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, _, err := shared.Setup(cmd.Flags())
		if err != nil {
			return err
		}

		groups := env.Groups
		if len(groups) == 0 {
			return fmt.Errorf("No groups found")
		}

		fmt.Println("The following groups was found:")
		fmt.Println()

		for _, group := range groups {
			fmt.Printf("  '%s' with alias '%s' (line %d to %d)", group, slug.Make(group.String()), group.Position.FirstLine, group.Position.LastLine)
			fmt.Println()
		}

		return nil
	},
}
