package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type colorType int

const (
	Danger colorType = 1 << iota
	Dark
	Info
	Light
	NoColor
	Primary
	Secondary
	Success
	Warning
)

type Writer struct {
	cache  map[colorType]Printer
	theme  Theme
	writer *lipgloss.Renderer
}

func (w Writer) Danger() Printer {
	return w.Color(Danger)
}

func (w Writer) Dark() Printer {
	return w.Color(Dark)
}

func (w Writer) Info() Printer {
	return w.Color(Info)
}

func (w Writer) Light() Printer {
	return w.Color(Light)
}

func (w Writer) NoColor() Printer {
	return w.Color(NoColor)
}

func (w Writer) Primary() Printer {
	return w.Color(Primary)
}

func (w Writer) Secondary() Printer {
	return w.Color(Secondary)
}

func (w Writer) Success() Printer {
	return w.Color(Success)
}

func (w Writer) Warning() Printer {
	return w.Color(Warning)
}

func (w Writer) Color(colorType colorType) Printer {
	if printer, ok := w.cache[colorType]; ok {
		return printer
	}

	var style Style

	switch colorType {
	case Danger:
		style = w.theme.Danger

	case Dark:
		style = w.theme.Dark

	case Info:
		style = w.theme.Info

	case Light:
		style = w.theme.Light

	case Primary:
		style = w.theme.Primary

	case Secondary:
		style = w.theme.Secondary

	case Success:
		style = w.theme.Success

	case Warning:
		style = w.theme.Warning

	case NoColor:
		style = w.theme.NoColor
	}

	w.cache[colorType] = style.NewPrinter(w.writer)

	return w.cache[colorType]
}
