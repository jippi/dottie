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

type ThemeWriter struct {
	cache  map[colorType]Printer
	theme  Theme
	writer *lipgloss.Renderer
}

func (tw ThemeWriter) Danger() Printer {
	return tw.Color(Danger)
}

func (tw ThemeWriter) Dark() Printer {
	return tw.Color(Dark)
}

func (tw ThemeWriter) Info() Printer {
	return tw.Color(Info)
}

func (tw ThemeWriter) Light() Printer {
	return tw.Color(Light)
}

func (tw ThemeWriter) NoColor() Printer {
	return tw.Color(NoColor)
}

func (tw ThemeWriter) Primary() Printer {
	return tw.Color(Primary)
}

func (tw ThemeWriter) Secondary() Printer {
	return tw.Color(Secondary)
}

func (tw ThemeWriter) Success() Printer {
	return tw.Color(Success)
}

func (tw ThemeWriter) Warning() Printer {
	return tw.Color(Warning)
}

func (tw ThemeWriter) Color(colorType colorType) Printer {
	if printer, ok := tw.cache[colorType]; ok {
		return printer
	}

	var style Style

	switch colorType {
	case Danger:
		style = tw.theme.Danger

	case Dark:
		style = tw.theme.Dark

	case Info:
		style = tw.theme.Info

	case Light:
		style = tw.theme.Light

	case Primary:
		style = tw.theme.Primary

	case Secondary:
		style = tw.theme.Secondary

	case Success:
		style = tw.theme.Success

	case Warning:
		style = tw.theme.Warning

	case NoColor:
		style = tw.theme.NoColor
	}

	tw.cache[colorType] = style.NewPrinter(tw.writer)

	return tw.cache[colorType]
}
