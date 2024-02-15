package tui

import "github.com/charmbracelet/lipgloss"

type Box struct {
	Header lipgloss.Style
	Body   lipgloss.Style
}

func (b Box) Copy() Box {
	return Box{
		Header: b.Header.Copy(),
		Body:   b.Body.Copy(),
	}
}
