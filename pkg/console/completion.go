package console

import (
	"github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
)

func (m *model) commandSuggestions(root *cobra.Command) []string {
	var suggestions []string

	if root.HasAvailableSubCommands() {
		for _, c := range root.Commands() {
			if c.Hidden {
				continue
			}

			suggestions = append(suggestions, c.Name())
		}
	}

	return suggestions
}

func (m *model) findCommand() bool {
	if m.input.Value() == "" {
		return false
	}

	args, err := shellquote.Split(m.input.Value())

	cmd, _, err := m.rootCommand.Find(args)
	if err != nil {
		return false
	}

	m.currentCommand = cmd
}
