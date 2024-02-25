package console

import (
	"strings"

	"github.com/elk-language/go-prompt"
	"github.com/ionoscloudsdk/comptplus"
	"github.com/jippi/dottie/pkg/tui"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "console",
		Short: "Interactive Terminal UI (TUI) console",
		Args:  cobra.ExactArgs(0),
		RunE:  runE,
	}

	return cmd
}

func runE(cmd *cobra.Command, _ []string) error {
	root := cmd.Root()

	advancedPrompt := &comptplus.CobraPrompt{
		RootCmd:                  root,
		AddDefaultExitCommand:    true,
		PersistFlagValues:        true,
		ShowHelpCommandAndFlags:  true,
		DisableCompletionCommand: true,
		GoPromptOptions: []prompt.Option{
			prompt.WithTitle("dottie"),
			prompt.WithPrefix("dottie: "),
			prompt.WithMaxSuggestion(10),
		},
		DynamicSuggestionsFunc: func(cmd *cobra.Command, _ string, document *prompt.Document) []prompt.Suggest {
			suggestions := []prompt.Suggest{}

			if cmd.ValidArgsFunction == nil {
				return suggestions
			}

			// We do not have access to the "args" being completed, so we pass in an empty slice
			// which in many cases will trigger the "full" list of suggestions
			arguments, _ := cmd.ValidArgsFunction(cmd, []string{}, document.GetWordBeforeCursor())
			for _, name := range arguments {
				// ValidArgsFunction() returns "name\tdescription" pairs
				// so we split by "\t" to get the individual parts
				parts := strings.SplitN(name, "\t", 2)

				// Build the default suggest
				suggestion := prompt.Suggest{
					Text: parts[0],
				}

				// If the autocomplete has a description, add it to the suggest
				if len(parts) == 2 {
					suggestion.Description = parts[1]
				}

				suggestions = append(suggestions, suggestion)
			}

			return suggestions
		},
		OnErrorFunc: func(err error) {
			tui.StderrFromContext(cmd.Context()).Danger().Println("Error:", err)
		},
	}

	advancedPrompt.RunContext(cmd.Context())

	return nil
}
