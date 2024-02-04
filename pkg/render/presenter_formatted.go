package render

import (
	"strings"
	"unicode"

	"github.com/jippi/dottie/pkg/ast"
)

type FormattedPresenter struct {
	wrapped Presenter
	output  Outputter
}

var _ Presenter = (*FormattedPresenter)(nil)

func (r *FormattedPresenter) Statement(stmt any, previous ast.Statement, settings Settings) string {
	switch val := stmt.(type) {
	// When formatting .env file we ignore all new-lines as we control those directly
	case *ast.Newline:
		return ""

	case *ast.Assignment:
		output := r.wrapped.Statement(stmt, previous, settings)
		if len(output) == 0 {
			return ""
		}

		buff := LineBuffer{}

		// Looks like current and previous Statement is both "Assignment"
		// which mean they might be too close in the document, so we will
		// attempt to inject some new-lines to give them some space
		if settings.ShowPretty && val.Is(previous) {
			// only allow cuddling of assignments if they both have no comments
			if val.HasComments() || assignmentHasComments(previous) {
				buff.Newline()
			}
		}

		return buff.
			Add(output).
			Get()

	// Offload all other statements to the wrapped presenter
	default:
		return r.wrapped.Statement(stmt, previous, settings)
	}
}

// NewFormattedPresenter creates a presenter that will format the output
// such as trimming newlines, ensuring spacing between items and os on
func NewFormattedPresenter(settings Settings) *FormattedPresenter {
	var output Outputter = Plain{}

	if settings.WithColors() {
		output = Colorized{}
	}

	presenter := &FormattedPresenter{
		output:  output,
		wrapped: NewDirectPresenter(settings),
	}

	return presenter
}

func (r *FormattedPresenter) Document(doc *ast.Document, settings Settings) string {
	out := &LineBuffer{}

	// Root statements
	root := r.Statement(doc.Statements, doc, settings)
	if len(root) > 0 {
		out.Add(root)
	}

	// Groups
	out.Add(r.Statement(doc.Groups, doc, settings))

	if settings.ShowPretty {
		out.Newline()
	}

	return strings.TrimLeftFunc(out.Get(), unicode.IsSpace)
}

func (r *FormattedPresenter) Group(group *ast.Group, settings Settings) string {
	if !group.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	rendered := r.Statement(group.Statements, nil, settings)
	if len(rendered) == 0 {
		return ""
	}

	res := &LineBuffer{}

	if settings.WithGroups() && len(rendered) > 0 {
		res.Add(r.output.Group(group, settings))
	}

	if settings.ShowPretty {
		res.Newline()
	}

	// Render the statements attached to the group
	rendered = strings.TrimSpace(rendered)
	res.Add(rendered)

	if settings.ShowPretty {
		res.Newline()
	}

	return res.Get()
}

func (r *FormattedPresenter) Assignment(a *ast.Assignment, settings Settings) string {
	if !settings.Match(a) || !a.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	res := &LineBuffer{}
	res.Add(r.Statement(a.Comments, a, settings))

	return r.output.Assignment(a, settings)
}

func (r *FormattedPresenter) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	if settings.WithComments() {
		return ""
	}

	var parent ast.Statement

	if isAssignmentComment {
		parent = &ast.Assignment{}
	}

	return r.wrapped.Statement(comment, parent, settings)
}

func (r *FormattedPresenter) SetOutput(output Outputter) {
	r.output = output
}
