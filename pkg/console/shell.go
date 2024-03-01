package console

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

type (
	errMsg error
)

type model struct {
	input          textinput.Model
	senderStyle    lipgloss.Style
	err            error
	rootCommand    *cobra.Command
	currentCommand *cobra.Command
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		commands []tea.Cmd
		inputCmd tea.Cmd
	)

	m.input, inputCmd = m.input.Update(msg)
	commands = append(commands, inputCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			commands = append(commands, tea.Printf(">> Dump: %s", spew.Sdump(m.input)))
			commands = append(commands, tea.Printf(">> Position: %d", m.input.Position()))
			commands = append(commands, tea.Printf(">> Value: %s", m.input.Value()))
			commands = append(commands, tea.Printf(">> Runes: %v", msg.Runes))
			commands = append(commands, tea.Printf("dottie: %s", m.input.Value()))

			m.input.Reset()

		default:
			commands = append(commands, tea.Printf(">> Runes: %v (%s) | %s", msg.Runes, string(msg.Runes), m.input.Value()))
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg

		return m, nil
	}

	return m, tea.Sequence(commands...)
}

func (m model) View() string {
	return m.input.View()
}
