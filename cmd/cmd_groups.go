package main

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/urfave/cli/v3"
)

var groupsCommand = &cli.Command{
	Name:   "groups",
	Usage:  "Print groups found in the .env file",
	Before: setup,
	Action: func(_ context.Context, _ *cli.Command) error {
		groups := env.Groups
		if len(groups) == 0 {
			return fmt.Errorf("No groups found")
		}

		fmt.Println("The following groups was found:")
		fmt.Println()

		for _, group := range groups {
			fmt.Printf("  '%s' with alias '%s' (line %d to %d)", group, slug.Make(group.String()), group.FirstLine, group.LastLine)
			fmt.Println()
		}

		return nil
	},
}
