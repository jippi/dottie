package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type colorType int

const (
	Danger colorType = iota
	Dark
	Info
	Light
	NoColor
	Primary
	Secondary
	Success
	Warning
)

type ThemeWriter struct {
	cache  map[colorType]Printer
	theme  Theme
	writer *lipgloss.Renderer
}

func (tp ThemeWriter) Color(colorType colorType) Printer {
	if printer, ok := tp.cache[colorType]; ok {
		return printer
	}

	var style Style

	switch colorType {
	case Danger:
		style = tp.theme.Danger

	case Dark:
		style = tp.theme.Dark

	case Info:
		style = tp.theme.Info

	case Light:
		style = tp.theme.Light

	case Primary:
		style = tp.theme.Primary

	case Secondary:
		style = tp.theme.Secondary

	case Success:
		style = tp.theme.Success

	case Warning:
		style = tp.theme.Warning

	case NoColor:
		style = tp.theme.NoColor
	}

	tp.cache[colorType] = style.NewPrinter(tp.writer)

	return tp.cache[colorType]
}
