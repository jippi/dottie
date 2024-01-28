package main

import (
	"context"

	"dotfedi/pkg"
	"dotfedi/pkg/ast"

	"github.com/urfave/cli/v3"
)

func setup(_ context.Context, cmd *cli.Command) error {
	var err error
	env, err = pkg.Load(cmd.String("file"))
	if err != nil {
		return err
	}

	settings = &ast.RenderSettings{
		FilterKeyPrefix:  cmd.Root().String("filter-key-prefix"),
		FilterGroup:      cmd.Root().String("filter-group"),
		IncludeCommented: cmd.Root().Bool("include-commented"),

		ShowPretty:     cmd.Root().Bool("pretty"),
		ShowBlankLines: cmd.Root().Bool("with-blank-lines"),
		ShowComments:   cmd.Root().Bool("with-comments"),
		ShowGroups:     cmd.Root().Bool("with-groups"),
	}

	return nil
}
