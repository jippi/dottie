package render

import "strings"

type LineBuffer struct {
	lines []string
}

func (lb *LineBuffer) Add(str string) *LineBuffer {
	if len(str) == 0 {
		return lb
	}

	if str == "\n" {
		str = ""
	}

	lb.lines = append(lb.lines, str)

	return lb
}

func (lb *LineBuffer) AddPrinted(str string) bool {
	if len(str) == 0 {
		return false
	}

	if str == "\n" {
		str = ""
	}

	lb.lines = append(lb.lines, str)

	return true
}

func (lb *LineBuffer) Get() string {
	return strings.Join(lb.lines, "\n")
}

func (lb *LineBuffer) Newline() {
	lb.lines = append(lb.lines, "")
}
