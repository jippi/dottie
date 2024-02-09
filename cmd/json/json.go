package json

import (
	"encoding/json"
	"fmt"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "json",
	Short: "Print as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, _, err := shared.Setup(cmd.Flags())
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))

		return nil
	},
}
