package fmt

import (
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "fmt",
		Short:   "Format a .env file",
		GroupID: "manipulate",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := cmd.Flag("file").Value.String()

			env, err := pkg.Load(filename)
			if err != nil {
				return err
			}

			if err := pkg.Save(cmd.Context(), filename, env); err != nil {
				return err
			}

			tui.ColorPrinterFromContext(cmd.Context(), tui.Stdout, tui.Success).Printfln("File [%s] was successfully formatted", filename)

			return nil
		},
	}
}
