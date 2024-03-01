package console

import (
	"github.com/spf13/cobra"
)

func (m *model) commandSuggestions(cmd *cobra.Command) []string {
	suggestions := []string{}

	if !cmd.HasAvailableSubCommands() {
		return suggestions
	}

	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}

		suggestions = append(suggestions, c.Name())
	}

	return suggestions
}

func (m *model) refreshSuggestions(cmd *cobra.Command) {
	suggestions := []string{}
	suggestions = append(suggestions, m.commandSuggestions(cmd)...)

	m.input.SetSuggestions(suggestions)
}

func (m *model) findCommand() {
	if m.input.Value() == "" {
		m.currentCommand = nil

		return
	}

	args := SafeSplitWords(m.input.Value())

	cmd, _, err := m.rootCommand.Find(JoinWords(args))
	if err != nil {
		m.currentCommand = nil

		return
	}

	if m.currentCommand == cmd {
		return
	}

	m.currentCommand = cmd
	m.refreshSuggestions(cmd)
}
