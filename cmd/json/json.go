package json

import (
	"encoding/json"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "json",
	Short:   "Print as JSON",
	GroupID: "output",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := cmd.Flag("file").Value.String()

		env, warn, err := pkg.Load(filename)
		if warn != nil {
			tui.Theme.Warning.StderrPrinter().Println(warn)
		}
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
