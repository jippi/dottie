package render

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

type PlainRenderer struct{}

func (r *PlainRenderer) Document(doc *ast.Document, settings Settings) string {
	var buf bytes.Buffer

	// Root statements
	root := r.Statements(doc.Statements, settings)
	if len(root) > 0 {
		buf.WriteString(root)
	}

	// Groups
	hasOutput := settings.WithGroups() && len(root) > 0

	for _, group := range doc.Groups {
		output := r.Group(group, settings)

		if hasOutput && len(output) > 0 {
			buf.WriteString("\n")
		}

		hasOutput = settings.WithGroups() && len(output) > 0

		buf.WriteString(output)
	}

	return strings.TrimLeftFunc(buf.String(), unicode.IsSpace)
}

func (r *PlainRenderer) Statements(statements []ast.Statement, settings Settings) string {
	var (
		buf     bytes.Buffer
		prev    ast.Statement
		printed bool
	)

	for _, stmt := range statements {
		switch val := stmt.(type) {
		case *ast.Group:
			panic("group should never happen in renderStatements")

		case *ast.Comment:
			printed = true

			buf.WriteString(r.Comment(val, settings, false))

		case *ast.Assignment:
			output := r.Assignment(val, settings)
			if len(output) == 0 {
				continue
			}

			// Looks like current and previous is both "Assignment"
			// which mean they are too close in the document, so we will
			// attempt to inject some new-lines to give them some space
			if settings.WithBlankLines() && val.Is(prev) {
				// only allow cuddling of assignments if they both have no comments
				if val.HasComments() || assignmentHasComments(prev) {
					buf.WriteString("\n")
				}
			}

			buf.WriteString(output)

			printed = true

		case *ast.Newline:
			output := r.Newline(val, settings)
			if len(output) == 0 {
				continue
			}

			// Don't print multiple newlines after each other
			if val.Is(prev) {
				continue
			}

			buf.WriteString(output)
		}

		prev = stmt
	}

	// If nothing "useful" was printed, don't bother outputting the groups buffer
	if !printed {
		return ""
	}

	str := buf.String()

	// Remove any duplicate newlines that might have crept into the output
	if settings.WithBlankLines() {
		str = strings.TrimRightFunc(str, unicode.IsSpace)
	}

	return "\n" + str
}

func (r *PlainRenderer) Group(group *ast.Group, settings Settings) string {
	if !group.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	var buf bytes.Buffer

	rendered := r.Statements(group.Statements, settings)
	if len(rendered) == 0 {
		return ""
	}

	if settings.WithGroups() && len(rendered) > 0 {
		if settings.WithColors() {
			out := tui.Theme.Info.Printer(tui.RendererWithTTY(&buf))
			out.Println("################################################################################")
			out.ApplyStyle(tui.Bold).Println(group.Name)
			out.Println("################################################################################")
			out.Println()
		} else {
			buf.WriteString("################################################################################")
			buf.WriteString("\n")

			buf.WriteString(group.Name)
			buf.WriteString("\n")

			buf.WriteString("################################################################################")
			buf.WriteString("\n")
			buf.WriteString("\n")
		}
	}

	// Render the statements attached to the group
	buf.WriteString(strings.TrimFunc(rendered, unicode.IsSpace))

	if settings.WithBlankLines() {
		return "\n" + buf.String()
	}

	return buf.String()
}

func (r *PlainRenderer) Assignment(a *ast.Assignment, settings Settings) string {
	if !settings.Match(a) || !a.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	var buf bytes.Buffer

	if settings.WithComments() {
		for _, c := range a.Comments {
			buf.WriteString(r.Comment(c, settings, true))
		}
	}

	if !a.Active {
		if settings.WithColors() {
			out := tui.Theme.Danger.Printer(tui.RendererWithTTY(&buf))
			out.Print("#")
		} else {
			buf.WriteString("#")
		}
	}

	val := a.Literal

	if settings.Interpolate {
		val = a.Interpolated
	}

	if settings.WithColors() {
		var buf bytes.Buffer

		tui.Theme.Primary.Printer(tui.RendererWithTTY(&buf)).Print(a.Name)
		tui.Theme.Dark.Printer(tui.RendererWithTTY(&buf)).Print("=")
		tui.Theme.Success.Printer(tui.RendererWithTTY(&buf)).Print(a.Quote)
		tui.Theme.Warning.Printer(tui.RendererWithTTY(&buf)).Print(val)
		tui.Theme.Success.Printer(tui.RendererWithTTY(&buf)).Print(a.Quote)

		return buf.String()
	}

	// panic(a.Quote)
	buf.WriteString(fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, val, a.Quote))
	buf.WriteString("\n")

	return buf.String()
}

func (r *PlainRenderer) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	if !settings.WithComments() || (!isAssignmentComment && !comment.BelongsToGroup(settings.FilterGroup)) {
		return ""
	}

	if !settings.WithColors() {
		return comment.Value + "\n"
	}

	var buf bytes.Buffer
	out := tui.Theme.Success.Printer(tui.RendererWithTTY(&buf))

	if comment.Annotation == nil {
		out.Println(comment.Value)

		return buf.String()
	}

	if comment.Annotation != nil {
		out.Print("# ")
		out.ApplyStyle(tui.Bold).Print("@", comment.Annotation.Key)
		out.Print(" ")
		out.Println(comment.Annotation.Value)
	}

	return buf.String()
}

func (r *PlainRenderer) Newline(newline *ast.Newline, settings Settings) string {
	if !settings.WithBlankLines() {
		return ""
	}

	return "\n"
}
