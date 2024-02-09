package shared

import (
	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/render"
	"github.com/spf13/pflag"
)

func Setup(flags *pflag.FlagSet) (*ast.Document, *render.Settings, error) {
	boolFlag := func(name string) bool {
		return BoolFlag(flags, name)
	}

	stringFlag := func(name string) string {
		return StringFlag(flags, name)
	}

	env, err := pkg.Load(stringFlag("file"))
	if err != nil {
		return nil, nil, err
	}

	settings := render.NewSettings(
		render.WithOutputType(render.Plain),
		render.WithBlankLines(boolFlag("with-blank-lines")),
		render.WithColors(boolFlag("colors")),
		render.WithComments(boolFlag("with-comments")),
		render.WithFilterGroup(stringFlag("group")),
		render.WithFilterKeyPrefix(stringFlag("key-prefix")),
		render.WithGroupBanners(boolFlag("with-groups")),
		render.WithIncludeDisabled(boolFlag("include-commented")),
	)

	if boolFlag("pretty") {
		settings.Apply(render.WithFormattedOutput(true))
	}

	return env, settings, nil
}
