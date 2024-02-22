package groups

import (
	"errors"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "groups",
		Short:   "Print groups found in the .env file",
		Args:    cobra.ExactArgs(0),
		GroupID: "output",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := cmd.Flag("file").Value.String()

			document, err := pkg.Load(cmd.Context(), filename)
			if err != nil {
				return err
			}

			groups := document.Groups
			if len(groups) == 0 {
				return errors.New("No groups found")
			}

			maxWidth := longesGroupName(groups)

			stdout := tui.StdoutFromContext(cmd.Context())
			primary := stdout.Primary()
			secondary := stdout.Secondary()

			stdout.Info().Box("Groups in " + filename)

			for _, group := range groups {
				primary.Printf("%-"+strconv.Itoa(maxWidth)+"s", slug.Make(group.String()))
				primary.Print(" ")
				secondary.Printfln("(%s:%d)", filename, group.Position.FirstLine)
			}

			return nil
		},
	}
}

func longesGroupName(groups []*ast.Group) int {
	length := 0

	for _, group := range groups {
		if len(group.Name) > length {
			length = len(group.Name)
		}
	}

	return length
}
