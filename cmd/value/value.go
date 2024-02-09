package value

import (
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "value",
	Usage:     "Print value of a env key if it exists",
	ArgsUsage: "KEY",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		env, _, err := shared.Setup(ctx, cmd)
		if err != nil {
			return err
		}

		key := cmd.Args().Get(0)
		if len(key) == 0 {
			return fmt.Errorf("Missing required argument: KEY")
		}

		existing := env.Get(key)
		if existing == nil {
			return fmt.Errorf("Key [%s] does not exists", key)
		}

		if !existing.Active && !cmd.Bool("include-commented") {
			return fmt.Errorf("Key [%s] exists, but is commented out - use [--include-commented] to include it", key)
		}

		fmt.Println(existing.Interpolated)

		return nil
	},
}
