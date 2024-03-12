package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/ui/layout"
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

func NewModel(ctx context.Context, document *ast.Document) model {
	// Initialize a global zone manager, so we don't have to pass around the manager
	// throughout components.
	zone.NewGlobal()

	return model{
		document: document,
		footer:   layout.FooterModel{},
		form: form{
			ctx:      ctx,
			document: document,
		},
		groups: group{
			id:    zone.NewPrefix(),
			title: "Groups",
			items: Map(document.Groups, func(g *ast.Group) groupItem {
				return groupItem{name: g.String()}
			}),
		},
	}
}

type model struct {
	height int
	width  int
	ready  bool

	document *ast.Document

	// Models
	groups group
	form   form
	footer layout.FooterModel
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.groups.Init())
	cmds = append(cmds, m.form.Init())
	cmds = append(cmds, m.footer.Init())

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.ready = true
	}

	return m.propagate(msg)
}

func (m *model) propagate(msg tea.Msg, commands ...tea.Cmd) (tea.Model, tea.Cmd) {
	// Footer
	{
		var cmd tea.Cmd
		m.footer, cmd = m.footer.Update(msg)
		commands = append(commands, cmd)
	}

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

	return m, tea.Sequence(commands...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	// layout
	content := lipgloss.JoinHorizontal(lipgloss.Top, m.groups.View(), m.form.View())
	output := lipgloss.JoinVertical(lipgloss.Left, content, m.footer.View())

	return zone.Scan(output)
}
