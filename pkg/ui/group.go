package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(subtle).
			MarginRight(2)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItemStyle = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDoneStyle = func(s string) string {
		return checkMark + lipgloss.NewStyle().Render(s)
	}
)

type changeGroupMsg struct {
	name string
}

func changeGroupCmd(name string) tea.Cmd {
	return func() tea.Msg {
		return changeGroupMsg{name: name}
	}
}

type groupItem struct {
	name   string
	active bool
}

type group struct {
	id     string
	height int
	width  int

	title string
	items []groupItem
}

func (m group) Init() tea.Cmd {
	return nil
}

func (m group) Update(msg tea.Msg) (group, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.MouseMsg:
		if msg.Button != tea.MouseButtonLeft {
			return m, nil
		}

		// If the zone wasn't targeted at all, abort
		if !zone.Get(m.id).InBounds(msg) {
			return m, nil
		}

		var cmd tea.Cmd

		// Check individual items if they were targeted
		for i, item := range m.items {
			m.items[i].active = zone.Get(m.id + item.name).InBounds(msg)

			if m.items[i].active {
				cmd = changeGroupCmd(item.name)
			}
		}

		return m, cmd

	case changeGroupMsg:
		for i, item := range m.items {
			m.items[i].active = item.name == msg.name
		}
	}

	return m, nil
}

func (m group) View() string {
	out := []string{listHeader(m.title)}

	for _, item := range m.items {
		if item.active {
			out = append(out, zone.Mark(m.id+item.name, listDoneStyle(item.name)))

			continue
		}

		out = append(out, zone.Mark(m.id+item.name, listItemStyle(item.name)))
	}

	return listStyle.Render(
		zone.Mark(
			m.id,
			lipgloss.JoinVertical(
				lipgloss.Left,
				out...,
			),
		),
	)
}
