package print_cmd

import (
	"fmt"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/render"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "print",
		Short: "Print environment variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			env, settings, err := shared.Setup(cmd.Flags())
			if err != nil {
				return err
			}

			settings.Apply(render.WithInterpolation(true))

			fmt.Println(render.NewRenderer(*settings).Statement(env).String())

			return nil
		},
	}

	cmd.Flags().Bool("pretty", false, "implies --colors --with-comments --with-blank-lines --with-groups")
	cmd.Flags().Bool("colors", true, "Enable/disable color output")
	cmd.Flags().Bool("with-comments", false, "Show comments")
	cmd.Flags().Bool("with-blank-lines", false, "Show blank lines between sections")
	cmd.Flags().Bool("with-groups", false, "Show group banners")
	cmd.Flags().Bool("include-commented", false, "Include commented KEY/VALUE pairs")
	cmd.Flags().String("key-prefix", "", "Filter by key prefix")
	cmd.Flags().String("group", "", "Filter by group name")

	return cmd
}
