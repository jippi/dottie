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
	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/validation"
	zone "github.com/lrstanley/bubblezone"
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
	id       string
	name     string
	document *ast.Document

	focusIndex int
	inputs     []tea.Model
	ctx        context.Context
	viewport   viewport.Model
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
			if str == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if str == "up" || str == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			offset := 0

			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds = append(cmds, m.inputs[i].(*huh.Input).Focus()) //nolint

					m.viewport.SetYOffset(offset)

					continue
				}

				newLines := strings.Split(m.inputs[i].View(), "\n")
				offset += len(newLines) + 1

				cmds = append(cmds, m.inputs[i].(*huh.Input).Blur()) //nolint
			}
		}

	case changeGroupMsg:
		m.name = msg.name

		inputs := []tea.Model{}

		for i, field := range m.document.GetGroup(m.name).Assignments() {
			input := huh.NewInput().
				Title(field.Name).
				Value(&field.Literal).
				Key(field.Name).
				Description(strings.TrimSpace(field.Documentation(true))).
				Validate(func(s string) error {
					err := validator.New().Var(s, field.ValidationRules())
					if err != nil {
						z := ast.NewError(field, err)

						return errors.New(validation.Explain(m.ctx, m.document, z, field, false, false))
					}

					return nil
				})

			if i == 0 {
				cmds = append(cmds, input.Focus())
			}

			inputs = append(inputs, input)
		}

		m.inputs = inputs
	}

	// Handle character input and blinking
	return m.propagate(msg, cmds...)
}

func (m *form) propagate(msg tea.Msg, commands ...tea.Cmd) (tea.Model, tea.Cmd) {
	for i := range m.inputs {
		var cmd tea.Cmd

		m.inputs[i], cmd = m.inputs[i].Update(msg)
		commands = append(commands, cmd)
	}

	var cmd tea.Cmd

	m.viewport, cmd = m.viewport.Update(msg)
	commands = append(commands, cmd)

	m.viewport.SetContent(m.render())

	return m, tea.Batch(commands...)
}

func (m form) render() string {
	var buf strings.Builder

	for i := range m.inputs {
		buf.WriteString(zone.Mark(fmt.Sprintf("%s-%d", m.id, i), m.inputs[i].View()))

		if i < len(m.inputs)-1 {
			buf.WriteString("\n\n")
		}
	}

	return buf.String()
}

func (m form) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	return m.viewport.View()
}
