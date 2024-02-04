package render

import (
	"bytes"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

type Outputter interface {
	Group(*ast.Group, Settings) string
	Assignment(*ast.Assignment, Settings) string
	Comment(*ast.Comment, Settings, bool) string
	Newline(*ast.Newline, Settings) string
}

var _ Outputter = (*Colorized)(nil)

type Colorized struct{}

func (c Colorized) Group(group *ast.Group, settings Settings) string {
	res := &LineBuffer{}

	var buf bytes.Buffer

	out := tui.Theme.Info.Printer(tui.RendererWithTTY(&buf))
	out.Println("################################################################################")
	out.ApplyStyle(tui.Bold).Println(group.Name)
	out.Print("################################################################################")

	res.Add(buf.String())

	return res.Get()
}

func (c Colorized) Assignment(a *ast.Assignment, settings Settings) string {
	var buf bytes.Buffer

	if !a.Active {
		tui.Theme.Danger.BuffPrinter(&buf).Print("#")
	}

	val := a.Literal

	if settings.Interpolate {
		val = a.Interpolated
	}

	tui.Theme.Primary.BuffPrinter(&buf).Print(a.Name)
	tui.Theme.Dark.BuffPrinter(&buf).Print("=")
	tui.Theme.Success.BuffPrinter(&buf).Print(a.Quote)
	tui.Theme.Warning.BuffPrinter(&buf).Print(val)
	tui.Theme.Success.BuffPrinter(&buf).Print(a.Quote)

	return (&LineBuffer{}).
		Add(buf.String()).
		Get()
}

func (r Colorized) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	var buf bytes.Buffer

	out := tui.Theme.Success.BuffPrinter(&buf)

	if comment.Annotation == nil {
		out.Print(comment.Value)

		return buf.String()
	}

	if comment.Annotation != nil {
		out.Print("# ")
		out.ApplyStyle(tui.Bold).Print("@", comment.Annotation.Key)
		out.Print(" ")
		out.Print(comment.Annotation.Value)
	}

	return buf.String()
}

func (r Colorized) Newline(newline *ast.Newline, settings Settings) string {
	if newline.Blank && !settings.WithBlankLines() {
		return ""
	}

	return "\n"
}
