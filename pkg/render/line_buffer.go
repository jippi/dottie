package render

import (
	"strings"
	"unicode"
)

const Newline = "\n"

type LineBuffer struct {
	lines []string
}

func NewLineBuffer() *LineBuffer {
	return &LineBuffer{}
}

func (lb *LineBuffer) Add(str string) *LineBuffer {
	lb.AddAndReturnPrinted(str)

	return lb
}

func (lb *LineBuffer) AddAndReturnPrinted(str string) bool {
	if len(str) == 0 {
		return false
	}

	if str == Newline {
		str = ""
	}

	lb.lines = append(lb.lines, str)

	return true
}

func (lb *LineBuffer) AddNewline() *LineBuffer {
	lb.lines = append(lb.lines, "")

	return lb
}

func (lb *LineBuffer) Get() string {
	return strings.Join(lb.lines, Newline)
}

func (lb *LineBuffer) GetWithEOF() string {
	return strings.TrimRightFunc(lb.Get(), unicode.IsSpace) + Newline
}
