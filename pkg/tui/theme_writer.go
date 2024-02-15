package tui

import "io"

type ThemeWriter struct {
	cache map[colorType]Printer
	theme Theme
	w     io.Writer
}

func (tp ThemeWriter) Color(colorType colorType) Printer {
	if printer, ok := tp.cache[colorType]; ok {
		return printer
	}

	var color Color

	switch colorType {
	case Danger:
		color = tp.theme.Danger
	case Dark:
		color = tp.theme.Dark
	case Info:
		color = tp.theme.Info
	case Light:
		color = tp.theme.Light
	case Primary:
		color = tp.theme.Primary
	case Secondary:
		color = tp.theme.Secondary
	case Success:
		color = tp.theme.Success
	case Warning:
		color = tp.theme.Warning
	case NoColor:
		color = tp.theme.NoColor
	}

	tp.cache[colorType] = color.BufferPrinter(tp.w)

	return tp.cache[colorType]
}
