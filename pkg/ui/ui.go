package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jippi/dottie/pkg/ast"
	zone "github.com/lrstanley/bubblezone"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()
)

type model struct {
	height int
	width  int

	document *ast.Document
	groups   group
	form     tea.Model
	ready    bool
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.groups.Init())
	cmds = append(cmds, m.form.Init())

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.ready = true

		m.height = msg.Height
		m.width = msg.Width
	}

	return m.propagate(msg)
}

func (m *model) propagate(msg tea.Msg, commands ...tea.Cmd) (tea.Model, tea.Cmd) {
	// Groups
	{
		var cmd tea.Cmd
		m.groups, cmd = m.groups.Update(msg)
		commands = append(commands, cmd)
	}

	// Form
	{
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width).Padding(1, 2, 1, 2)

	return zone.Scan(
		s.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.groups.View(),
				m.form.View(),
			),
		),
	)
}
