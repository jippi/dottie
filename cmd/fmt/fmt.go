package fmt

import (
	"context"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "fmt",
	Usage: "Format the file",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		env, _, err := shared.Setup(ctx, cmd)
		if err != nil {
			return err
		}

		return pkg.Save(cmd.String("file"), env)
	},
}
