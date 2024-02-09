package json

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "json",
	Usage: "Print as JSON",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		env, _, err := shared.Setup(ctx, cmd)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))

		return nil
	},
}
