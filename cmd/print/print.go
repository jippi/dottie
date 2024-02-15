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
		Args:    cobra.NoArgs,
		GroupID: "output",
		RunE:    runE,
	}

	cmd.Flags().Bool("pretty", false, "implies --color --comments --blank-lines --group-banners")

	cmd.Flags().String("key-prefix", "", "Filter by key prefix")
	cmd.Flags().String("group", "", "Filter by group name")

	shared.BoolWithInverse(cmd, "blank-lines", true, "Show blank lines", "Do not show blank lines")
	shared.BoolWithInverse(cmd, "color", true, "Enable color output", "Disable color output")
	shared.BoolWithInverse(cmd, "commented", false, "Show disabled assignments", "Do not show disabled assignments")
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

	doc, err := pkg.Load(stringFlag("file"))
	if err != nil {
		return nil, nil, err
	}

	settings := render.NewSettings(
		render.WithOutputType(render.Plain),
		render.WithColors(shared.BoolWithInverseValue(flags, "color")),
		render.WithComments(shared.BoolWithInverseValue(flags, "comments")),
		render.WithBlankLines(shared.BoolWithInverseValue(flags, "blank-lines")),
		render.WithGroupBanners(shared.BoolWithInverseValue(flags, "group-banners")),
		render.WithIncludeDisabled(shared.BoolWithInverseValue(flags, "commented")),
		render.WithInterpolation(shared.BoolWithInverseValue(flags, "interpolation")),
		render.WithFilterGroup(stringFlag("group")),
		render.WithFilterKeyPrefix(stringFlag("key-prefix")),
	)

	var allErrors, allWarnings error

	if settings.InterpolatedValues {
		var warn, err error

		for _, assignment := range doc.AllAssignments() {
			if !assignment.Enabled {
				continue
			}

			warn, err = doc.InterpolateStatement(assignment)

			allWarnings = multierr.Append(allWarnings, warn)
			allErrors = multierr.Append(allErrors, err)
		}
	}

	if boolFlag("pretty") {
		settings.Apply(render.WithFormattedOutput(true))
	}

	tui.MaybePrintWarnings(cmd.Context(), allWarnings)

	return doc, settings, allErrors
}
