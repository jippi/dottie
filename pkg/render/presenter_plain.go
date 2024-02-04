package render

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

type PlainPresenter struct{}

var _ Presenter = (*PlainPresenter)(nil)

func (r *PlainPresenter) Statement(stmt any, previous ast.Statement, settings Settings) string {
	switch val := stmt.(type) {
	case *ast.Document:
		return r.Document(val, settings)

	case *ast.Group:
		return r.Group(val, settings)

	case *ast.Comment:
		return r.Comment(val, settings, false)

	case *ast.Assignment:
		return r.Assignment(val, settings)

	case *ast.Newline:
		return r.Newline(val, settings)

	default:
		panic(fmt.Sprintf("Unknown statement: %T", val))
	}
}

func (r *PlainPresenter) Statements(statements []ast.Statement, settings Settings) string {
	var (
		res     = &Accumulator{}
		prev    ast.Statement
		printed bool
	)

	for _, stmt := range statements {
		switch val := stmt.(type) {
		case *ast.Comment:
			printed = res.AddPrinted(r.Comment(val, settings, false))

		case *ast.Assignment:
			output := r.Assignment(val, settings)
			if len(output) == 0 {
				continue
			}

			// Looks like current and previous Statement is both "Assignment"
			// which mean they might be too close in the document, so we will
			// attempt to inject some new-lines to give them some space
			if settings.WithBlankLines() && val.Is(prev) {
				// only allow cuddling of assignments if they both have no comments
				if val.HasComments() || assignmentHasComments(prev) {
					res.Newline()
				}
			}

			printed = res.AddPrinted(output)

		case *ast.Newline:
			if settings.WithBlankLines() {
				continue
			}

			output := r.Newline(val, settings)
			if len(output) == 0 {
				continue
			}

			// Don't print multiple newlines after each other
			if val.Is(prev) {
				continue
			}

			if prev != nil && !assignmentHasComments(prev) {
				continue
			}

			res.Newline()

		default:
			panic(fmt.Sprintf("Unknown statement: %T", val))
		}

		prev = stmt
	}

	// If nothing "useful" was printed, don't bother outputting the groups buffer
	if !printed {
		return ""
	}

	str := res.Get()

	// Remove any duplicate newlines that might have crept into the output
	if settings.WithBlankLines() {
		str = strings.TrimRightFunc(str, unicode.IsSpace)
	}

	return str
}

func (r *PlainPresenter) Document(doc *ast.Document, settings Settings) string {
	out := &Accumulator{}

	// Root statements
	root := r.Statements(doc.Statements, settings)
	if len(root) > 0 {
		out.Add(root)
	}

	// Groups
	hasOutput := settings.WithGroups() && len(root) > 0

	for _, group := range doc.Groups {
		output := r.Group(group, settings)

		if hasOutput && len(output) > 0 {
			out.Newline()
		}

		hasOutput = settings.WithGroups() && len(output) > 0

		if settings.WithBlankLines() {
			output = strings.TrimSpace(output)
		}

		out.Add(output)
	}

	if settings.WithBlankLines() {
		out.Newline()
	}

	return strings.TrimLeftFunc(out.Get(), unicode.IsSpace)
}

func (r *PlainPresenter) Group(group *ast.Group, settings Settings) string {
	if !group.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	rendered := r.Statements(group.Statements, settings)
	if len(rendered) == 0 {
		return ""
	}

	res := &Accumulator{}

	if settings.WithGroups() && len(rendered) > 0 {
		if settings.WithColors() {
			var buf bytes.Buffer

			out := tui.Theme.Info.Printer(tui.RendererWithTTY(&buf))
			out.Println("################################################################################")
			out.ApplyStyle(tui.Bold).Println(group.Name)
			out.Print("################################################################################")

			res.Add(buf.String())
		} else {
			res.Add("################################################################################")
			res.Add(group.Name)
			res.Add("################################################################################")
		}
	}

	if settings.WithBlankLines() {
		res.Newline()
	}

	// Render the statements attached to the group
	rendered = strings.TrimSpace(rendered)
	res.Add(rendered)

	if settings.WithBlankLines() {
		res.Newline()
	}

	return res.Get()
}

func (r *PlainPresenter) Assignment(a *ast.Assignment, settings Settings) string {
	if !settings.Match(a) || !a.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	res := &Accumulator{}

	if settings.WithComments() {
		for _, c := range a.Comments {
			res.Add(r.Comment(c, settings, true))
		}
	}

	var buf bytes.Buffer

	val := a.Literal

	if settings.Interpolate {
		val = a.Interpolated
	}

	if !a.Active {
		if settings.WithColors() {
			tui.Theme.Danger.BuffPrinter(&buf).Print("#")
		} else {
			buf.WriteString("#")
		}
	}

	if settings.WithColors() {
		tui.Theme.Primary.BuffPrinter(&buf).Print(a.Name)
		tui.Theme.Dark.BuffPrinter(&buf).Print("=")
		tui.Theme.Success.BuffPrinter(&buf).Print(a.Quote)
		tui.Theme.Warning.BuffPrinter(&buf).Print(val)
		tui.Theme.Success.BuffPrinter(&buf).Print(a.Quote)

		return res.
			Add(buf.String()).
			Get()
	}

	return res.
		Add(fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, val, a.Quote)).
		Get()
}

func (r *PlainPresenter) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	if !settings.WithComments() || (!isAssignmentComment && !comment.BelongsToGroup(settings.FilterGroup)) {
		return ""
	}

	if !settings.WithColors() {
		return comment.Value
	}

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

func (r *PlainPresenter) Newline(newline *ast.Newline, settings Settings) string {
	if !settings.WithBlankLines() {
		return ""
	}

	return "\n"
}
