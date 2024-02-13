package fmt

import (
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "fmt",
	Short:   "Format a .env file",
	GroupID: "manipulate",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := cmd.Flag("file").Value.String()

		env, warn, err := pkg.Load(filename)
		if warn != nil {
			tui.Theme.Warning.StderrPrinter().Println(warn)
		}
		if err != nil {
			return err
		}

		if err := pkg.Save(filename, env); err != nil {
			return err
		}

		tui.Theme.Success.StdoutPrinter().Printfln("File [%s] was successfully formatted", filename)

		return nil
	},
}
