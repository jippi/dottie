package main

import (
	"context"
	"fmt"

	"dotfedi/pkg"
	"dotfedi/pkg/ast"

	"github.com/urfave/cli/v3"
)

var setCommand = &cli.Command{
	Name:      "set",
	Usage:     "Set a key/value pair",
	Before:    setup,
	ArgsUsage: "KEY VALUE",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "commented",
			OnlyOnce: true,
		},
		&cli.BoolFlag{
			Name:     "error-if-missing",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "group",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "before",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "after",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "quote-style",
			Usage:    "single|double|none",
			OnlyOnce: true,
		},
		&cli.StringSliceFlag{
			Name: "comment",
		},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		key := cmd.Args().Get(0)
		if len(key) == 0 {
			return fmt.Errorf("Missing required argument: KEY")
		}

		var group *ast.Group
		value := cmd.Args().Get(1)

		assignment := env.Get(key)
		if assignment == nil {
			if cmd.Bool("error-if-missing") {
				return fmt.Errorf("Key [%s] does not exists", key)
			}

			group = env.GetGroup(ast.RenderSettings{FilterGroup: cmd.String("group")})
			if group == nil {
				group = &ast.Group{Name: cmd.String("group")}
				env.Groups = append(env.Groups, group)
			}

			assignment = &ast.Assignment{
				Key:   key,
				Group: group,
			}

			switch {
			case len(cmd.String("before")) > 0:
				before := cmd.String("before")

				var res []ast.Statement
				for _, stmt := range group.Statements {
					x, ok := stmt.(*ast.Assignment)
					if !ok {
						res = append(res, stmt)
						continue
					}

					if x.Key == before {
						res = append(res, assignment)
					}

					res = append(res, stmt)
				}

				group.Statements = res

			default:
				group.Statements = append(group.Statements, assignment)
			}
		}

		env.Assignments = append(env.Assignments, assignment)

		assignment.Value = value
		assignment.Commented = cmd.Bool("commented")
		assignment.SetQuote(cmd.String("quote-style"))

		if comments := cmd.StringSlice("comment"); len(comments) > 0 {
			slice := make([]*ast.Comment, 0)

			for _, v := range comments {
				slice = append(slice, ast.NewComment(v))
			}

			assignment.Comments = slice
		}

		return pkg.Save(cmd.String("file"), env)
	},
}
