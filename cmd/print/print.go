package print_cmd

import (
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "print",
		Short:   "Print environment variables",
		Args:    cobra.ExactArgs(0),
		GroupID: "output",
		RunE:    runE,
	}

	cmd.Flags().Bool("pretty", false, "implies --color --comments --blank-lines --group-banners")
	cmd.Flags().Bool("export", false, "prefix all key/value pairs with [export] statement")
	cmd.Flags().Bool("with-disabled", false, "Include disabled assignments")

	cmd.Flags().String("key-prefix", "", "Filter by key prefix")
	cmd.Flags().String("group", "", "Filter by group name")

	shared.BoolWithInverse(cmd, "blank-lines", true, "Show blank lines", "Do not show blank lines")
	shared.BoolWithInverse(cmd, "color", true, "Enable color output", "Disable color output")
	shared.BoolWithInverse(cmd, "comments", false, "Show comments", "Do not show comments")
	shared.BoolWithInverse(cmd, "group-banners", false, "Show group banners", "Do not show group banners")
	shared.BoolWithInverse(cmd, "interpolation", true, "Enable interpolation", "Disable interpolation")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	document, settings, err := setup(cmd)
	if err != nil {
		return err
	}

	tui.StdoutFromContext(cmd.Context()).
		NoColor().
		Println(
			render.NewRenderer(*settings).
				Statement(cmd.Context(), document).
				String(),
		)

	return nil
}

func setup(cmd *cobra.Command) (*ast.Document, *render.Settings, error) {
	flags := cmd.Flags()

	boolFlag := func(name string) bool {
		return shared.BoolFlag(flags, name)
	}

	stringFlag := func(name string) string {
		return shared.StringFlag(flags, name)
	}

	doc, err := pkg.Load(cmd.Context(), stringFlag("file"))
	if err != nil {
		return nil, nil, err
	}

	settings := render.NewSettings(
		render.WithBlankLines(shared.BoolWithInverseValue(flags, "blank-lines")),
		render.WithColors(shared.BoolWithInverseValue(flags, "color")),
		render.WithComments(shared.BoolWithInverseValue(flags, "comments")),
		render.WithFilterGroup(stringFlag("group")),
		render.WithFilterKeyPrefix(stringFlag("key-prefix")),
		render.WithGroupBanners(shared.BoolWithInverseValue(flags, "group-banners")),
		render.WithIncludeDisabled(boolFlag("with-disabled")),
		render.WithInterpolation(shared.BoolWithInverseValue(flags, "interpolation")),
		render.WithOutputType(render.Plain),
	)

	var allErrors error

	if settings.InterpolatedValues {
		var err error

		for _, assignment := range doc.AllAssignments() {
			if !assignment.Enabled {
				continue
			}

			err = doc.InterpolateStatement(cmd.Context(), assignment, false)

			allErrors = multierr.Append(allErrors, err)
		}
	}

	if boolFlag("pretty") {
		settings.Apply(render.WithFormattedOutput(true))
	}

	if boolFlag("export") {
		settings.Apply(render.WithExport(true))
	}

	return doc, settings, allErrors
}
