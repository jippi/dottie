package render

import (
	"bytes"
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

var _ Outputter = (*Plain)(nil)

type Plain struct{}

func (c Plain) Group(group *ast.Group, settings Settings) string {
	out := NewLineBuffer()

	out.Add("################################################################################")
	out.Add(group.Name)
	out.Add("################################################################################")

	return out.Get()
}

func (c Plain) Assignment(a *ast.Assignment, settings Settings) string {
	var buf bytes.Buffer

	if !a.Active {
		buf.WriteString("#")
	}

	val := a.Literal

	if settings.Interpolate {
		val = a.Interpolated
	}

	buf.WriteString(fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, val, a.Quote))

	return buf.String()
}

func (r Plain) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	return comment.Value
}

func (r Plain) Newline(newline *ast.Newline, settings Settings) string {
	if newline.Blank && !settings.WithBlankLines() {
		return ""
	}

	return "\n"
}
