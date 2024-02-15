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
		Args:    cobra.ExactArgs(0),
		GroupID: "output",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := cmd.Flag("file").Value.String()

			document, err := pkg.Load(filename)
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(document, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			return nil
		},
	}
}
