package render

import (
	"os"
	"strings"
	"unicode"
)

const Newline = "\n"

type LineBufferItem struct {
	Literal string
	Debug   string
}

type LineBuffer struct {
	lines []LineBufferItem
}

func NewLineBuffer() *LineBuffer {
	return &LineBuffer{}
}

func (lb *LineBuffer) Empty() bool {
	return lb == nil || len(lb.lines) == 0
}

func (lb *LineBuffer) AddString(in string) *LineBuffer {
	if len(in) == 0 {
		return lb
	}

	if in == Newline {
		return lb.AddNewline("LineBuffer:AddString")
	}

	lb.lines = append(lb.lines, LineBufferItem{
		Literal: in,
	})

	return lb
}

func (lb *LineBuffer) Add(buf *LineBuffer) *LineBuffer {
	lb.AddAndReturnPrinted(buf)

	return lb
}

func (lb *LineBuffer) AddAndReturnPrinted(buf *LineBuffer) bool {
	if buf == nil || len(buf.lines) == 0 {
		return false
	}

	lb.lines = append(lb.lines, buf.lines...)

	return true
}

func (lb *LineBuffer) AddNewline(id ...string) *LineBuffer {
	str := "# AddNewline"

	if len(id) > 0 {
		str = "# " + strings.Join(id, " ")
	}

	lb.lines = append(lb.lines, LineBufferItem{Literal: "", Debug: str})

	return lb
}

func (lb *LineBuffer) Get() string {
	res := []string{}

	for _, l := range lb.lines {
		if os.Getenv("DEBUG") == "1" && l.Debug != "" {
			res = append(res, l.Debug)

			continue
		}

		res = append(res, l.Literal)
	}

	return strings.Join(res, Newline)
}

func (lb *LineBuffer) Trim() *LineBuffer {
	var start, stop int

	for start = 0; start < len(lb.lines); start++ {
		if lb.lines[start].Literal != "" {
			break
		}
	}

	for stop = len(lb.lines); stop >= 0; {
		if lb.lines[stop-1].Literal != "" {
			break
		}

		stop--
	}

	lb.lines = lb.lines[start:stop]

	return lb
}

func (lb *LineBuffer) GetWithEOF() string {
	return strings.TrimRightFunc(lb.Get(), unicode.IsSpace) + Newline
}

func (lb *LineBuffer) GetTrimmed() string {
	return strings.TrimRightFunc(lb.Get(), unicode.IsSpace)
}
