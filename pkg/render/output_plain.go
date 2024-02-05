package render

import (
	"bytes"
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

var _ Output = (*PlainOutput)(nil)

type PlainOutput struct{}

func (PlainOutput) GroupBanner(group *ast.Group, settings Settings) *LineBuffer {
	out := NewLineBuffer()

	out.AddString("################################################################################")
	out.AddString(group.Name)
	out.AddString("################################################################################")

	return out
}

func (PlainOutput) Assignment(a *ast.Assignment, settings Settings) *LineBuffer {
	var buf bytes.Buffer

	if !a.Active {
		buf.WriteString("#")
	}

	val := a.Literal

	if settings.UseInterpolatedValues {
		val = a.Interpolated
	}

	buf.WriteString(fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, val, a.Quote))

	return NewLineBuffer().AddString(buf.String())
}

func (r PlainOutput) Comment(comment *ast.Comment, settings Settings) *LineBuffer {
	return NewLineBuffer().AddString(comment.Value)
}

func (r PlainOutput) Newline(newline *ast.Newline, settings Settings) *LineBuffer {
	if newline.Blank && !settings.WithBlankLines() {
		return nil
	}

	return NewLineBuffer().AddNewline("PlainOutput:Newline")
}
