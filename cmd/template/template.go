package template

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Short:   "Render a template",
		Args:    cobra.MinimumNArgs(1),
		GroupID: "output",
		RunE:    runE,
	}

	shared.BoolWithInverse(cmd, "interpolation", true, "Enable interpolation", "Disable interpolation")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	document, err := setup(cmd)
	if err != nil {
		return err
	}

	return template.Must(template.New(cmd.Flags().Arg(0)).Funcs(sprig.FuncMap()).ParseFiles(cmd.Flags().Args()...)).Execute(cmd.OutOrStdout(), document)
}

func setup(cmd *cobra.Command) (*ast.Document, error) {
	flags := cmd.Flags()

	stringFlag := func(name string) string {
		return shared.StringFlag(flags, name)
	}

	boolFlag := func(name string) bool {
		return shared.BoolFlag(flags, name)
	}

	doc, err := pkg.Load(cmd.Context(), stringFlag("file"))
	if err != nil {
		return nil, err
	}

	var allErrors error

	if boolFlag("interpolation") {
		for _, assignment := range doc.AllAssignments() {
			err := doc.InterpolateStatement(cmd.Context(), assignment, boolFlag("with-disabled"))

			allErrors = multierr.Append(allErrors, err)
		}
	}

	return doc, allErrors
}
