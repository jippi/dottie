package enable

import (
	"context"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "enable",
	Usage:     "Uncomment/enable a key if it exists",
	ArgsUsage: "KEY",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		env, _, err := shared.Setup(ctx, cmd)
		if err != nil {
			return err
		}

		key := cmd.Args().Get(0)

		existing := env.Get(key)
		existing.Active = true

		return pkg.Save(cmd.String("file"), env)
	},
}
