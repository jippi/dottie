package ui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/validation"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type form struct {
	id        string
	groupName string
	document  *ast.Document

	focusIndex int
	ctx        context.Context
	viewport   viewport.Model
	fields     []tea.Model
	errors     map[string]string
	ready      bool
}

func (m form) Init() tea.Cmd {
	return changeGroupCmd(m.document.Groups[0].String())
}

func (m form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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
					key.WithKeys("u", "ctrl+u"),
					key.WithHelp("u", "½ page up"),
				),
				HalfPageDown: key.NewBinding(
					key.WithKeys("d", "ctrl+d"),
					key.WithHelp("d", "½ page down"),
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
			m.errors = make(map[string]string)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
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

			if m.focusIndex > len(m.fields) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.fields)
			}

			offset := 0

			for i := 0; i <= len(m.fields)-1; i++ {
				if i == m.focusIndex {
					cmds = append(cmds, m.fields[i].(*huh.Input).Focus()) //nolint

					m.viewport.SetYOffset(offset)

					continue
				}

				newLines := strings.Split(m.fields[i].View(), "\n")
				offset += len(newLines) + 1

				cmds = append(cmds, m.fields[i].(*huh.Input).Blur()) //nolint
			}
		}

	case changeGroupMsg:
		m.groupName = msg.name

		fields := []tea.Model{}

		for i, field := range m.document.GetGroup(m.groupName).Assignments() {
			field := field

			input := huh.NewInput().
				Title(field.Name).
				Value(&field.Literal).
				Key(field.Name).
				Description(strings.TrimSpace(field.Documentation(true))).
				Validate(func(s string) error {
					if m.errors == nil {
						m.errors = make(map[string]string)
					}

					valErr, err := m.document.ValidateSingleAssignment(m.ctx, field, nil, nil)
					if err != nil {
						m.errors[field.Name] = err.Error()

						return err
					}

					if valErr != nil {
						m.errors[field.Name] = validation.Explain(m.ctx, m.document, valErr, field, false, false)

						return errors.New(m.errors[field.Name])
					}

					delete(m.errors, field.Name)

					return nil
				})

			if i == 0 {
				cmds = append(cmds, input.Focus())
			}

			fields = append(fields, input)
		}

		m.fields = fields
	}

	// Handle character input and blinking
	return m.propagate(msg, cmds...)
}

func (m *form) propagate(msg tea.Msg, commands ...tea.Cmd) (tea.Model, tea.Cmd) {
	for i, field := range m.fields {
		var cmd tea.Cmd

		m.fields[i], cmd = field.Update(msg)
		commands = append(commands, cmd)
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	commands = append(commands, cmd)

	var output string

	for _, field := range m.fields {
		output += field.View() + "\n\n"

		realField := field.(*huh.Input) //nolint

		if err := m.errors[realField.GetKey()]; len(err) > 0 {
			output += "Error: " + err + "\n"
		}
	}

	m.viewport.SetContent(output)

	return m, tea.Batch(commands...)
}

func (m form) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	return m.viewport.View()
}
