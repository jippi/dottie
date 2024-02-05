package render

import (
	"bytes"
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

var _ Output = (*PlainOutput)(nil)

type PlainOutput struct{}

func (PlainOutput) GroupBanner(group *ast.Group, settings Settings) *Lines {
	out := NewLinesCollection()

	out.Add("################################################################################")
	out.Add(group.Name)
	out.Add("################################################################################")

	return out
}

func (PlainOutput) Assignment(a *ast.Assignment, settings Settings) *Lines {
	var buf bytes.Buffer

	if !a.Active {
		buf.WriteString("#")
	}

	val := a.Literal

	if settings.useInterpolatedValues {
		val = a.Interpolated
	}

	buf.WriteString(fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, val, a.Quote))

	return NewLinesCollection().Add(buf.String())
}

func (r PlainOutput) Comment(comment *ast.Comment, settings Settings) *Lines {
	return NewLinesCollection().Add(comment.Value)
}

func (r PlainOutput) Newline(newline *ast.Newline, settings Settings) *Lines {
	if newline.Blank && !settings.ShowBlankLines() {
		return nil
	}

	return NewLinesCollection().Newline("PlainOutput:Newline")
}
