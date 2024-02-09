package enable

import (
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "enable",
	Usage:     "Enable (uncomment) a KEY if it exists",
	ArgsUsage: "KEY",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		key := cmd.Args().Get(0)
		if len(key) == 0 {
			return fmt.Errorf("Missing required argument: KEY")
		}

		env, _, err := shared.Setup(ctx, cmd)
		if err != nil {
			return err
		}

		existing := env.Get(key)
		if existing == nil {
			return fmt.Errorf("Could not find KEY [%s]", key)
		}

		existing.Enable()

		return pkg.Save(cmd.String("file"), env)
	},
}
