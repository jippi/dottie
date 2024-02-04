package render

import "strings"

type Accumulator struct {
	lines []string
}

func (a *Accumulator) Add(str string) *Accumulator {
	if len(str) == 0 {
		return a
	}

	a.lines = append(a.lines, str)

	return a
}

func (a *Accumulator) AddPrinted(str string) bool {
	if len(str) == 0 {
		return false
	}

	a.lines = append(a.lines, str)

	return true
}

func (a *Accumulator) Get() string {
	return strings.Join(a.lines, "\n")
}

func (a *Accumulator) Newline() {
	a.lines = append(a.lines, "")
}
