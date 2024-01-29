package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

var jsonCommand = &cli.Command{
	Name:   "json",
	Usage:  "Print as JSON",
	Before: setup,
	Action: func(_ context.Context, cmd *cli.Command) error {
		b, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))
		return nil
	},
}
