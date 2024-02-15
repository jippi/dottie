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
		Args:    cobra.NoArgs,
		GroupID: "manipulate",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := cmd.Flag("file").Value.String()

			document, err := pkg.Load(filename)
			if err != nil {
				return err
			}

			if err := pkg.Save(cmd.Context(), filename, document); err != nil {
				return err
			}

			tui.StdoutFromContext(cmd.Context()).
				Success().
				Printfln("File [ %s ] was successfully formatted", filename)

			return nil
		},
	}
}
