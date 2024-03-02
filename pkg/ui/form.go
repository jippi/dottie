package ui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/ui/component/textinput"
	"github.com/jippi/dottie/pkg/validation"
	zone "github.com/lrstanley/bubblezone"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()
)

type form struct {
	id        string
	groupName string
	document  *ast.Document
	theme     tui.Writer

	focusIndex int
	ctx        context.Context
	viewport   viewport.Model
	fields     []textinput.Model
	ready      bool
}

func (m form) Init() tea.Cmd {
	return changeGroupCmd(m.document.Groups[0].String())
}

func (m form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	m.theme = tui.StderrFromContext(m.ctx)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.KeyMap = viewport.KeyMap{
				PageDown: key.NewBinding(
					key.WithKeys("pgdown"),
					key.WithHelp("f/pgdn", "page down"),
				),
				PageUp: key.NewBinding(
					key.WithKeys("pgup"),
					key.WithHelp("b/pgup", "page up"),
				),
				HalfPageUp: key.NewBinding(
					key.WithKeys("ctrl+u"),
				),
				HalfPageDown: key.NewBinding(
					key.WithKeys("ctrl+d"),
				),
				Up: key.NewBinding(
					key.WithKeys("up"),
					key.WithHelp("↑/k", "up"),
				),
				Down: key.NewBinding(
					key.WithKeys("down"),
					key.WithHelp("↓/j", "down"),
				),
			}
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft {
			// Check individual items if they were targeted
			for idx := range m.fields {
				if zone.Get(fmt.Sprintf("%s-%d", m.id, idx)).InBounds(msg) {
					m.fields[idx].Focus()
				} else {
					m.fields[idx].Blur()
				}
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			str := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if str == "enter" && m.focusIndex == len(m.fields) {
				return m, tea.Quit
			}

			// Cycle indexes
			if str == "up" || str == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex >= len(m.fields) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.fields) - 1
			}

			offset := 0
			height := m.viewport.Height

			for i := 0; i <= len(m.fields)-1; i++ {
				if i == m.focusIndex {
					commands = append(commands, m.fields[i].Focus()) //nolint

					if offset > height {
						m.viewport.SetYOffset(offset)
					}

					continue
				}

				newLines := strings.Split(m.fields[i].View(), "\n")
				offset += len(newLines) + 1

				m.fields[i].Blur()
			}
		}

	case changeGroupMsg:
		m.groupName = msg.name

		var fields []textinput.Model

		for i, assignment := range m.document.GetGroup(m.groupName).Assignments() {
			assignment := assignment

			input := textinput.New(assignment)
			input.Prompt = m.promptFor(assignment, input)
			input.SetValue(assignment.Literal)
			input.Width = m.viewport.Width
			input.Validate = func(s string) error {
				assignment.Literal = s

				valErr, err := m.document.ValidateSingleAssignment(m.ctx, assignment, nil, nil)
				if err != nil {
					return err
				}

				if valErr != nil {
					return errors.New(validation.Explain(m.ctx, m.document, valErr, assignment, false, false))
				}

				return nil
			}

			if i == 0 {
				commands = append(commands, input.Focus())
			}

			fields = append(fields, input)
		}

		m.fields = fields
	}

	// Handle character input and blinking
	return m.propagate(msg, commands...)
}

func (m *form) propagate(msg tea.Msg, commands ...tea.Cmd) (tea.Model, tea.Cmd) {
	for i := range m.fields {
		var cmd tea.Cmd
		m.fields[i], cmd = m.fields[i].Update(msg)
		commands = append(commands, cmd)

		m.fields[i].Prompt = m.promptFor(m.fields[i].Assignment, m.fields[i])
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	m.viewport.SetContent(m.renderFields())

	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func (m form) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	return m.viewport.View()
}

func (m *form) renderFields() string {
	var output string

	for idx, field := range m.fields {
		str := field.View()

		if err := field.Err; err != nil {
			str += "\n\n" + lipgloss.NewStyle().Foreground(tui.Red500).Render(strings.TrimSpace(err.Error()))
		}

		str = zone.Mark(fmt.Sprintf("%s-%d", m.id, idx), strings.TrimSpace(str))

		output += lipgloss.NewStyle().Padding(1).Render(strings.TrimSpace(str))
	}

	return output
}

func (m *form) promptFor(assignment *ast.Assignment, field textinput.Model) string {
	docs := m.theme.Success().Sprint(assignment.Documentation(false))
	name := m.theme.Primary().Sprint(assignment.Name)

	if field.Err != nil {
		name = m.theme.Danger().Sprint(assignment.Name)
	}

	return strings.TrimSpace(docs) + strings.TrimSpace(name) + m.theme.Dark().Sprint("=")
}
