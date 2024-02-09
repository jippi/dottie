package set

import (
	"fmt"
	"strings"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set KEY=VALUE [KEY=VALUE ...]",
		Short: "Set/update one or multiple key=value pairs",
		RunE: func(cmd *cobra.Command, args []string) error {
			env, _, err := shared.Setup(cmd.Flags())
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return fmt.Errorf("Missing required argument: KEY=VALUE")
			}

			comments, _ := cmd.Flags().GetStringArray("comment")
			options := ast.UpsertOptions{
				InsertBefore:   shared.StringFlag(cmd.Flags(), "before"),
				Comments:       comments,
				ErrorIfMissing: shared.BoolFlag(cmd.Flags(), "error-if-missing"),
				Group:          shared.StringFlag(cmd.Flags(), "group"),
				SkipValidation: !shared.BoolFlag(cmd.Flags(), "validate"),
			}

			for _, stringPair := range args {
				pairSlice := strings.SplitN(stringPair, "=", 2)
				if len(pairSlice) != 2 {
					return fmt.Errorf("expected KEY=VALUE pair, missing '='")
				}

				key := pairSlice[0]
				value := pairSlice[1]

				assignment := &ast.Assignment{
					Name:    key,
					Literal: value,
					// by default we take the user input and assume its interpolated,
					// it will be interpolated inside (*Document).Set if applicable
					Interpolated: value,
					Active:       !shared.BoolFlag(cmd.Flags(), "disabled"),
					Quote:        token.QuoteFromString(shared.StringFlag(cmd.Flags(), "quote-style")),
				}

				//
				// Upsert key
				//

				assignment, err := env.Upsert(assignment, options)
				if err != nil {
					validation.Explain(env, validation.NewError(assignment, err))

					return fmt.Errorf("failed to upsert the key/value pair [%s]", key)
				}

				tui.Theme.Success.StderrPrinter().Printfln("Key [%s] was successfully upserted", key)
			}

			//
			// Save file
			//

			if err := pkg.Save(shared.StringFlag(cmd.Flags(), "file"), env); err != nil {
				return fmt.Errorf("failed to save file: %w", err)
			}

			tui.Theme.Success.StderrPrinter().Println("File was successfully saved")

			return nil
		},
	}

	cmd.Flags().Bool("disabled", false, "Set/change the flag to be disabled (commented out)")
	cmd.Flags().Bool("validate", true, "Validate the VALUE input before saving the file")
	cmd.Flags().Bool("error-if-missing", false, "Exit with an error if the KEY does not exists in the .env file already")
	cmd.Flags().String("group", "", "The (optional) group name to add the KEY=VALUE pair under")
	cmd.Flags().String("before", "", "If the key doesn't exist, add it to the file *before* this KEY")
	cmd.Flags().String("after", "", "If the key doesn't exist, add it to the file *after* this KEY")
	cmd.Flags().String("quote-style", "double", "The quote style to use (single, double, none)")
	cmd.Flags().StringSlice("comment", nil, "Set one or multiple lines of comments to the KEY=VALUE pair")

	return cmd
}
