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

var Command = &cobra.Command{
	Use:     "groups",
	Short:   "Print groups found in the .env file",
	GroupID: "output",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := cmd.Flag("file").Value.String()

		env, err := pkg.Load(filename)
		if err != nil {
			return err
		}

		groups := env.Groups
		if len(groups) == 0 {
			return errors.New("No groups found")
		}

		width := longesGroupName(groups)

		light := tui.Theme.Secondary.StdoutPrinter()
		key := tui.Theme.Primary.StdoutPrinter()
		info := tui.Theme.Info.StdoutPrinter()
		info.Box("Groups in " + filename)

		for _, group := range groups {
			key.Printf("%-"+strconv.Itoa(width)+"s", slug.Make(group.String()))
			key.Print(" ")
			light.Printfln("(%s:%d)", filename, group.Position.FirstLine)
		}

		return nil
	},
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
