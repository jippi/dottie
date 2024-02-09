package groups

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "groups",
	Usage: "Print groups found in the .env file",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		env, _, err := shared.Setup(ctx, cmd)
		if err != nil {
			return err
		}

		groups := env.Groups
		if len(groups) == 0 {
			return fmt.Errorf("No groups found")
		}

		fmt.Println("The following groups was found:")
		fmt.Println()

		for _, group := range groups {
			fmt.Printf("  '%s' with alias '%s' (line %d to %d)", group, slug.Make(group.String()), group.Position.FirstLine, group.Position.LastLine)
			fmt.Println()
		}

		return nil
	},
}
