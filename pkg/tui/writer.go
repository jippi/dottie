package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type style int

const (
	Danger style = 1 << iota
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
	cache  map[style]Printer
	theme  Theme
	writer *lipgloss.Renderer
}

func (w Writer) Danger() Printer {
	return w.Style(Danger)
}

func (w Writer) Dark() Printer {
	return w.Style(Dark)
}

func (w Writer) Info() Printer {
	return w.Style(Info)
}

func (w Writer) Light() Printer {
	return w.Style(Light)
}

func (w Writer) NoColor() Printer {
	return w.Style(NoColor)
}

func (w Writer) Primary() Printer {
	return w.Style(Primary)
}

func (w Writer) Secondary() Printer {
	return w.Style(Secondary)
}

func (w Writer) Success() Printer {
	return w.Style(Success)
}

func (w Writer) Warning() Printer {
	return w.Style(Warning)
}

func (w Writer) Style(colorType style) Printer {
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
