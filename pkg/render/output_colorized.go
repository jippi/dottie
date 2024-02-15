package render

import (
	"bytes"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

var _ Output = (*ColorizedOutput)(nil)

type ColorizedOutput struct{}

func (ColorizedOutput) GroupBanner(group *ast.Group, settings Settings) *Lines {
	var buf bytes.Buffer

	out := tui.Theme.Info.Printer(tui.RendererWithTTY(&buf))

	out.Println("################################################################################")
	out.ApplyStyle(tui.Bold).Println(group.Name)
	out.Print("################################################################################")

	return NewLinesCollection().Add(buf.String())
}

func (ColorizedOutput) Assignment(assignment *ast.Assignment, settings Settings) *Lines {
	var buf bytes.Buffer

	if !assignment.Enabled {
		tui.Theme.Danger.Printer(tui.RendererWithTTY(&buf)).Print("#")
	}

	val := assignment.Literal

	if settings.InterpolatedValues {
		val = assignment.Interpolated
	}

	tui.Theme.Primary.Printer(tui.RendererWithTTY(&buf)).Print(assignment.Name)
	tui.Theme.Dark.Printer(tui.RendererWithTTY(&buf)).Print("=")
	tui.Theme.Success.Printer(tui.RendererWithTTY(&buf)).Print(assignment.Quote)
	tui.Theme.Warning.Printer(tui.RendererWithTTY(&buf)).Print(val)
	tui.Theme.Success.Printer(tui.RendererWithTTY(&buf)).Print(assignment.Quote)

	return NewLinesCollection().Add(buf.String())
}

func (ColorizedOutput) Comment(comment *ast.Comment, settings Settings) *Lines {
	var buf bytes.Buffer

	out := tui.Theme.Success.Printer(tui.RendererWithTTY(&buf))

	if comment.Annotation == nil {
		out.Print(comment.Value)

		return NewLinesCollection().Add(buf.String())
	}

	if comment.Annotation != nil {
		out.Print("# ")
		out.ApplyStyle(tui.Bold).Print("@", comment.Annotation.Key)
		out.Print(" ")
		out.Print(comment.Annotation.Value)
	}

	return NewLinesCollection().Add(buf.String())
}

func (ColorizedOutput) Newline(newline *ast.Newline, settings Settings) *Lines {
	if newline.Blank && !settings.ShowBlankLines() {
		return nil
	}

	return NewLinesCollection().Newline("ColorizedOutput:Newline")
}
