package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type Writer struct {
	cache    map[styleIdentifier]Printer
	theme    Theme
	renderer *lipgloss.Renderer
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

func (w Writer) Style(colorType styleIdentifier) Printer {
	if printer, ok := w.cache[colorType]; ok {
		return printer
	}

	w.cache[colorType] = w.theme.Style(colorType).NewPrinter(w.renderer)

	return w.cache[colorType]
}
