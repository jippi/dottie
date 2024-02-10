package print_cmd

import (
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print",
		Short: "Print environment variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			env, settings, err := setup(cmd.Flags())
			if err != nil {
				return err
			}

			fmt.Println(render.NewRenderer(*settings).Statement(env).String())

			return nil
		},
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

func setup(flags *pflag.FlagSet) (*ast.Document, *render.Settings, error) {
	boolFlag := func(name string) bool {
		return shared.BoolFlag(flags, name)
	}

	stringFlag := func(name string) string {
		return shared.StringFlag(flags, name)
	}

	env, err := pkg.Load(stringFlag("file"))
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

	if boolFlag("pretty") {
		settings.Apply(render.WithFormattedOutput(true))
	}

	return env, settings, nil
}
