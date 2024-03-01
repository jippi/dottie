package console

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func NewModel(cmd *cobra.Command) model {
	root := cmd.Root()

	input := textinput.New()
	input.Placeholder = ""
	input.Prompt = "dottie: "
	input.ShowSuggestions = true
	input.Focus()

	return model{
		input:       input,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		rootCommand: root,
	}
}
