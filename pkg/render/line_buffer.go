package render

import (
	"strings"
	"unicode"
)

const Newline = "\n"

type LineBuffer struct {
	lines []string
}

func (lb *LineBuffer) Add(str string) *LineBuffer {
	lb.AddPrinted(str)

	return lb
}

func (lb *LineBuffer) AddPrinted(str string) bool {
	if len(str) == 0 {
		return false
	}

	if str == Newline {
		str = ""
	}

	lb.lines = append(lb.lines, str)

	return true
}

func (lb *LineBuffer) Get() string {
	return strings.Join(lb.lines, Newline)
}

func (lb *LineBuffer) GetWithEOF() string {
	return strings.TrimRightFunc(lb.Get(), unicode.IsSpace) + Newline
}

func (lb *LineBuffer) Newline() *LineBuffer {
	lb.lines = append(lb.lines, "")

	return lb
}

func (lb *LineBuffer) EnsureEOF() *LineBuffer {
	idx := len(lb.lines) - 1
	if idx > 0 && lb.lines[idx] != "" {
		return lb.Newline()
	}

	return lb
}
