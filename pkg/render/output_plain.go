package render

import (
	"bytes"
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

var _ Output = (*PlainOutput)(nil)

type PlainOutput struct{}

func (PlainOutput) GroupBanner(group *ast.Group, settings Settings) string {
	out := NewLineBuffer()

	out.Add("################################################################################")
	out.Add(group.Name)
	out.Add("################################################################################")

	return out.Get()
}

func (PlainOutput) Assignment(a *ast.Assignment, settings Settings) string {
	var buf bytes.Buffer

	if !a.Active {
		buf.WriteString("#")
	}

	val := a.Literal

	if settings.UseInterpolatedValues {
		val = a.Interpolated
	}

	buf.WriteString(fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, val, a.Quote))

	return buf.String()
}

func (r PlainOutput) Comment(comment *ast.Comment, settings Settings) string {
	return comment.Value
}

func (r PlainOutput) Newline(newline *ast.Newline, settings Settings) string {
	if newline.Blank && !settings.WithBlankLines() {
		return ""
	}

	return "\n"
}
