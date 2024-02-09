package render

import (
	"os"
	"strings"
)

const Newline = "\n"

// Line represents a line in an output
type Line struct {
	Literal string // The *real* value of the line being added
	Debug   string // A debug value to help track where NewLines are coming from (use DEBUG=1 to see them when using --pretty)
}

func (l Line) String() string {
	if os.Getenv("DEBUG") == "1" && l.Debug != "" {
		return l.Debug
	}

	return l.Literal
}

// Lines is a collection of lines
type Lines struct {
	lines []Line
}

func NewLinesCollection() *Lines {
	return &Lines{}
}

// IsEmpty returns if the Line has any lines or not
func (lb *Lines) IsEmpty() bool {
	return lb == nil || len(lb.lines) == 0
}

// Add a new arbitrary string to the Lines instance
func (lb *Lines) Add(value string) *Lines {
	// If there is no length, discard input
	if len(value) == 0 {
		return lb
	}

	// If the input is a single newline, redirect
	// the call to the dedicated Newline func
	// and track origin of the newline
	if value == Newline {
		return lb.Newline("LineBuffer:AddString")
	}

	// Append the Line to the collection
	lb.lines = append(lb.lines, Line{Literal: value})

	return lb
}

// Append adds one Lines collection to the end of another
// if its's not empty
func (lb *Lines) Append(buf *Lines) *Lines {
	if buf == nil || buf.IsEmpty() {
		return lb
	}

	lb.lines = append(lb.lines, buf.lines...)

	return lb
}

// Newline adds a Newline to the collection
//
// Optionally accepts a list of strings to help identify
// what the "origin" or reason for the newline is.
//
// This can be seen when using [DEBUG=1] environment variable
// along with [--pretty] CLI flag
func (lb *Lines) Newline(id ...string) *Lines {
	str := "# Lines.Newline"

	if len(id) > 0 {
		str = "# " + strings.Join(id, " ")
	}

	lb.lines = append(lb.lines, Line{Literal: "", Debug: str})

	return lb
}

// String joins all the LineItems together into slice and
// Join them by a newline + an additional trailing newline.
func (lb Lines) String() string {
	return strings.Join(lb.Lines(), Newline) + Newline
}

// Lines returns the raw slice of lines
func (lb Lines) Lines() []string {
	res := []string{}

	for _, line := range lb.lines {
		res = append(res, line.String())
	}

	return res
}
