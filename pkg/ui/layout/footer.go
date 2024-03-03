package layout

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))
)

var footerStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(subtle).
	MarginRight(1).
	MarginRight(0).
	BorderBottom(false).
	BorderTop(true)

type ShowHiddenMsg struct {
	Hide bool
}

func ShowHiddenCmd(val bool) tea.Cmd {
	return func() tea.Msg {
		return ShowHiddenMsg{
			Hide: val,
		}
	}
}

type FooterModel struct {
	height     int
	width      int
	ready      bool
	hideHidden bool
}

func (m FooterModel) Init() tea.Cmd {
	return nil
}

func (m FooterModel) Update(msg tea.Msg) (FooterModel, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.ready = true

	case tea.MouseMsg:
		if msg.String() != "left press" {
			break
		}

		if !zone.Get("footer/toggle-hidden").InBounds(msg) {
			break
		}

		m.hideHidden = !m.hideHidden

		return m, ShowHiddenCmd(m.hideHidden)
	}

	return m, tea.Batch(commands...)
}

func (m FooterModel) View() string {
	if !m.ready {
		return "Initializing ..."
	}

	statusKey := statusStyle.Render("SHOWING HIDDEN")
	if m.hideHidden {
		statusKey = statusStyle.Render("HIDING HIDDEN")
	}

	encoding := encodingStyle.Render(time.Now().Format(time.StampMilli))
	fishCake := fishCakeStyle.Render("🍥 Fish Cake")
	statusVal := statusText.Copy().
		Width(m.width - lipgloss.Width(statusKey) - lipgloss.Width(encoding) - lipgloss.Width(fishCake)).
		Render("Ravishing")

	bar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		zone.Mark("footer/toggle-hidden", statusKey),
		statusVal,
		encoding,
		fishCake,
	)

	return bar
}
