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
	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/ui/component/textinput"
	"github.com/jippi/dottie/pkg/ui/layout"
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
	selectors []ast.Selector

	focusIndex int
	ctx        context.Context
	viewport   viewport.Model
	fields     []textinput.Model
	ready      bool
}

func (m form) Init() tea.Cmd {
	return changeGroupCmd(m.document.Groups[0].String())
}

func (m form) Update(msg tea.Msg) (form, tea.Cmd) {
	var commands []tea.Cmd

	m.theme = tui.StderrFromContext(m.ctx)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-FooterHeight)
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
			m.viewport.Height = msg.Height - FooterHeight
		}

	case tea.MouseMsg:
		if msg.String() != "left press" {
			break
		}

		// Check individual items if they were targeted
		for idx := range m.fields {
			if zone.Get(fmt.Sprintf("%s-%d", m.id, idx)).InBounds(msg) {
				m.focusIndex = idx
				m.fields[idx].Focus()
			} else {
				m.fields[idx].Blur()
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			if !m.showingCompletion() {
				str := msg.String()

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
		}

	case layout.ShowHiddenMsg:
		m.selectors = []ast.Selector{}

		if msg.Hide {
			m.selectors = []ast.Selector{ast.ExcludeDisabledAssignments}
		}

		commands = append(commands, m.populateFields())

	case changeGroupMsg:
		m.focusIndex = 0
		m.groupName = msg.name

		commands = append(commands, m.populateFields())
	}

	// Handle character input and blinking
	return m.propagate(msg, commands...)
}

func (m *form) propagate(msg tea.Msg, commands ...tea.Cmd) (form, tea.Cmd) {
	// Fields
	{
		for i := range m.fields {
			var cmd tea.Cmd
			m.fields[i], cmd = m.fields[i].Update(msg)
			commands = append(commands, cmd)

			m.fields[i].Prompt = m.promptFor(m.fields[i].Assignment, m.fields[i])
		}
	}

	// Viewport
	{
		// Don't update viewport if we're doing field completion (to avoid hijacking keyboard input)
		if !m.showingCompletion() {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			commands = append(commands, cmd)
		}

		m.viewport.SetContent(
			listHeader.
				Width(m.viewport.Width).
				Render("Group: "+m.groupName) + m.renderFields(),
		)
	}

	return *m, tea.Batch(commands...)
}

func (m *form) showingCompletion() bool {
	if len(m.fields) == 0 {
		return false
	}

	return m.fields[m.focusIndex].ShowingAcceptSuggestion()
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
	prefix := ""

	if field.Err != nil {
		name = m.theme.Danger().Sprint(assignment.Name)
	}

	if !assignment.Enabled {
		docs = m.theme.Dark().Sprint(assignment.Documentation(false))
		name = m.theme.Dark().Sprint(assignment.Name)
		prefix = m.theme.Dark().Sprint("#")
	}

	return strings.TrimSpace(docs) + prefix + strings.TrimSpace(name) + m.theme.Dark().Sprint("=")
}

func (m *form) populateFields() tea.Cmd {
	var (
		commands []tea.Cmd
		fields   = make([]textinput.Model, 0)
	)

	for i, assignment := range m.document.GetGroup(m.groupName).Assignments(m.selectors...) {
		assignment := assignment

		input := textinput.New(assignment)
		input.ShowSuggestions = true

		for rule := parseRules(assignment.ValidationRules()); rule != nil; rule = rule.Next {
			if rule.Tag == "oneof" {
				input.SetSuggestions(strings.Split(rule.Param, " "))
			}
		}

		input.Prompt = m.promptFor(assignment, input)
		input.PromptStyle = lipgloss.NewStyle()
		input.TextStyle = lipgloss.NewStyle()
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

	return tea.Batch(commands...)
}

type cTag struct {
	Tag      string
	Param    string
	HasParam bool
	Next     *cTag
}

const (
	tagSeparator    = ","
	orSeparator     = "|"
	tagKeySeparator = "="
	utf8HexComma    = "0x2C"
	utf8Pipe        = "0x7C"
)

func parseRules(tag string) *cTag {
	if len(tag) == 0 {
		return nil
	}

	var (
		alias   string
		first   *cTag
		current *cTag
		tags    = strings.Split(tag, tagSeparator)
	)

	for i := 0; i < len(tags); i++ {
		alias = tags[i]

		if i == 0 {
			current = &cTag{}
			first = current
		} else {
			current.Next = &cTag{}
			current = current.Next
		}

		current.Tag = tags[i]

		// if a pipe character is needed within the param you must use the utf8Pipe representation "0x7C"
		orGroups := strings.Split(alias, orSeparator)

		for groupIdx := 0; groupIdx < len(orGroups); groupIdx++ {
			name, params, _ := strings.Cut(orGroups[groupIdx], tagKeySeparator)

			current.Tag = name

			if groupIdx > 0 {
				current.Next = &cTag{}
				current = current.Next
			}

			current.HasParam = len(params) > 0

			if len(current.Tag) == 0 {
				panic("invalid 1: " + tag + " || " + alias + " || " + name + " || " + params + " || " + orGroups[groupIdx] + " || " + spew.Sdump(current) + " || " + spew.Sdump(first))
			}

			if len(params) > 1 {
				current.Param = strings.Replace(strings.Replace(params, utf8HexComma, ",", -1), utf8Pipe, "|", -1)
			}
		}
	}

	return first
}
