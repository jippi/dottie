package main

import (
	"context"
	"fmt"

	"dotfedi/pkg/ast"

	"github.com/urfave/cli/v3"
)

var printCommand = &cli.Command{
	Name:   "print",
	Usage:  "Print environment variables",
	Before: setup,
	Action: func(_ context.Context, cmd *cli.Command) error {
		for _, s := range env.Statements {
			switch v := s.(type) {
			case *ast.Group:
				if len(filters.Group) > 0 && v.Comment != filters.Group {
					continue
				}

				if cmd.Bool("pretty") || cmd.Bool("with-groups") {
					fmt.Println("################################################################################")
					fmt.Println("# " + v.Comment)
					fmt.Println("################################################################################")
				}

			case *ast.Comment:
				if len(filters.Group) > 0 && (v.Group == nil || v.Group.Comment != filters.Group) {
					continue
				}

				if cmd.Bool("pretty") || cmd.Bool("with-comments") {
					fmt.Println(v)

					continue
				}

			case *ast.Assignment:
				if !filters.Match(v) {
					continue
				}

				if cmd.Bool("pretty") || cmd.Bool("with-comments") {
					fmt.Println(v)

					continue
				}

				fmt.Println(v.Assignment())

			case *ast.Newline:
				if len(filters.Group) > 0 && (v.Group == nil || v.Group.Comment != filters.Group) {
					continue
				}

				if cmd.Bool("pretty") || cmd.Bool("with-blank-lines") && v.Blank {
					fmt.Println()
				}
			}
		}

		return nil
	},
}
