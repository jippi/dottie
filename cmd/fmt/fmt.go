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
		Args:    cobra.ExactArgs(0),
		GroupID: "manipulate",
		RunE:    runE,
	}
}

func runE(cmd *cobra.Command, args []string) error {
	filename := cmd.Flag("file").Value.String()

	document, err := pkg.Load(cmd.Context(), filename)
	if err != nil {
		return err
	}

	if err := pkg.Save(cmd.Context(), filename, document); err != nil {
		return err
	}

	tui.StdoutFromContext(cmd.Context()).
		Success().
		Printfln("File was successfully formatted")

	return nil
}
