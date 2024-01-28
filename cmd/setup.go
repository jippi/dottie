package main

import (
	"context"

	"dotfedi/pkg"
	"dotfedi/pkg/filter"

	"github.com/urfave/cli/v3"
)

func setup(_ context.Context, cmd *cli.Command) error {
	var err error
	env, err = pkg.Load(cmd.String("file"))
	if err != nil {
		return err
	}

	filters = &filter.Filter{
		KeyPrefix: cmd.Root().String("filter-key-prefix"),
		Group:     cmd.Root().String("filter-group"),
	}

	return nil
}
