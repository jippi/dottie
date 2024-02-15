package json

import (
	"encoding/json"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "json",
		Short:   "Print as JSON",
		Args:    cobra.NoArgs,
		GroupID: "output",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := cmd.Flag("file").Value.String()

			env, err := pkg.Load(filename)
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
}
