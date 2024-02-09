package print_cmd

import (
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "print",
	Usage: "Print environment variables",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		env, settings, err := shared.Setup(ctx, cmd)
		if err != nil {
			return err
		}

		settings.Apply(render.WithInterpolation(true))

		fmt.Println(render.NewRenderer(*settings).Statement(env).String())

		return nil
	},
}
