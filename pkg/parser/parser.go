package parser

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/token"
	slogctx "github.com/veqryn/slog-context"
)

// Scanner converts a sequence of characters into a sequence of tokens.
type Scanner interface {
	NextToken(ctx context.Context) token.Token
}

// Parser takes a Scanner and builds an abstract syntax tree.
type Parser struct {
	filename      string
	scanner       Scanner
	token         token.Token
	previousToken token.Token
}

// New returns new Parser.
func New(ctx context.Context, scanner Scanner, filename string) *Parser {
	// If we're under test, the filename will be completely random, which will
	// mess with golden file output, so we override the filename to a static
	// value to make the tests deterministic.
	if testing.Testing() && strings.HasPrefix(filename, "/") {
		filename = "/fake/testing/path/" + filepath.Base(filename)
	}

	ctx = slogctx.With(
		ctx,
		slog.Group(
			"parser_state",
			slog.String("filename", filename),
		),
	)

	return &Parser{
		filename: filename,
		scanner:  scanner,
		token:    scanner.NextToken(ctx),
	}
}

// Parse parses the .env file and returns an ast.Statement.
func (p *Parser) Parse(ctx context.Context) (document *ast.Document, err error) {
	var (
		comments          []*ast.Comment
		currentGroup      *ast.Group
		previousStatement ast.Statement
		statementIndex    int
	)

	document = ast.NewDocument()

	for p.token.Type != token.EOF {
		ctx := slogctx.With(
			ctx,
			slog.Group("parser_state",
				slog.String("filename", p.filename),
				slog.Any("token", p.token),
			),
			slog.String("source", "parser"),
		)

		stmt, err := p.parseStatement(ctx)
		if err != nil {
			return nil, err
		}

		slogctx.Debug(ctx, "Parser.Parse() processing statement", slog.Any("statement", stmt))

		switch val := stmt.(type) {
		case *ast.Group:
			val.Position.File = p.filename

			// Track the last line of this group
			if currentGroup != nil {
				currentGroup.Position.LastLine = p.token.LineNumber
			}

			// Change the current group
			currentGroup = val

			// Append the group
			document.Groups = append(document.Groups, currentGroup)

			previousStatement = val

		case *ast.Assignment:
			val.Position.File = p.filename
			val.Position.Index = statementIndex

			statementIndex++

			// Assign accumulated comments to this assignment
			val.Comments = comments

			if len(val.Comments) > 0 {
				val.Position.FirstLine = val.Comments[0].Position.Line
			} else {
				val.Position.FirstLine = val.Position.Line
			}

			val.Position.LastLine = val.Position.Line

			// Assign the assignment to a grouping if such exists
			if currentGroup != nil {
				val.Group = currentGroup
				currentGroup.Statements = append(currentGroup.Statements, val)
			} else {
				document.Statements = append(document.Statements, stmt)
			}

			// Reset comment block
			comments = make([]*ast.Comment, 0, 1)
			previousStatement = val

		case *ast.Comment:
			val.Position.File = p.filename

			if val.Annotation != nil {
				document.Annotations = append(document.Annotations, val)
			}

			comments = append(comments, val)
			previousStatement = val

		case *ast.Newline:
			if !val.Blank {
				continue
			}

			val.Position.File = p.filename

			// If the previous statement was an assignment, ignore the newline
			// as we will be emitted that ourself later
			if val.Is(previousStatement) {
				last, _ := previousStatement.(*ast.Newline)
				last.Position.LastLine = val.Position.Line
				last.Repeated++

				continue
			}

			// If there is a blank line, print all previous comments
			for _, comment := range comments {
				if currentGroup != nil {
					comment.Group = currentGroup

					currentGroup.Statements = append(currentGroup.Statements, comment)
				} else {
					document.Statements = append(document.Statements, comment)
				}
			}

			// Attach the newline to a group for easier filtering
			if currentGroup != nil {
				val.Group = currentGroup
			}

			if currentGroup != nil {
				currentGroup.Statements = append(currentGroup.Statements, val)
			} else {
				document.Statements = append(document.Statements, val)
			}

			// Reset the accumulated comments slice
			comments = make([]*ast.Comment, 0, 1)
			previousStatement = val
		}
	}

	if currentGroup != nil {
		currentGroup.Position.LastLine = p.token.LineNumber
	}

	if len(comments) > 0 {
		if currentGroup != nil {
			for _, c := range comments {
				currentGroup.Statements = append(currentGroup.Statements, c)
			}
		} else {
			for _, c := range comments {
				document.Statements = append(document.Statements, c)
			}
		}
	}

	// Make sure to initialize the document so dependencies and such are
	// computed immediately
	document.Initialize(ctx)

	return document, nil
}

func (p *Parser) parseStatement(ctx context.Context) (ast.Statement, error) {
	slogctx.Debug(ctx, "Parser.parseStatement()")

	switch p.token.Type {
	case token.Identifier:
		return p.parseRowStatement(ctx)

	case token.Comment, token.CommentAnnotation:
		return p.parseCommentStatement(ctx)

	case token.GroupBanner:
		return p.parseGroupStatement(ctx)

	case token.EOF:
		return nil, nil

	case token.NewLine:
		res := &ast.Newline{
			Blank: p.wasEmptyLine(),
			Position: ast.Position{
				Line:      p.token.LineNumber,
				FirstLine: p.token.LineNumber,
				LastLine:  p.token.LineNumber,
			},
		}

		p.nextToken(ctx)

		return res, nil

	default:
		return nil, fmt.Errorf("(B) unexpected statement: %s(%q)", p.token.Type, p.token.Literal)
	}
}

func (p *Parser) parseGroupStatement(ctx context.Context) (ast.Statement, error) {
	slogctx.Debug(ctx, "Parser.parseGroupStatement()")

	var group *ast.Group

	p.nextToken(ctx)
	p.skipBlankLine(ctx)

	switch p.token.Type {
	case token.Comment:
		group = &ast.Group{
			Name: p.token.Literal,
			Position: ast.Position{
				FirstLine: p.token.LineNumber,
				Line:      p.token.LineNumber,
				LastLine:  p.token.LineNumber,
			},
		}

	default:
		return p.unexpectedToken("parseGroupStatement 1")
	}

	p.nextToken(ctx)
	p.skipBlankLine(ctx)

	switch p.token.Type {
	case token.GroupBanner:
		p.nextToken(ctx)

		return group, nil

	default:
		return p.unexpectedToken("parseGroupStatement 2")
	}
}

func (p *Parser) parseCommentStatement(ctx context.Context) (ast.Statement, error) {
	slogctx.Debug(ctx, "Parser.parseCommentStatement()")

	stm := &ast.Comment{
		Value:      p.token.Literal,
		Annotation: p.token.Annotation,
		Position: ast.Position{
			Line:      p.token.LineNumber,
			FirstLine: p.token.LineNumber,
			LastLine:  p.token.LineNumber,
		},
	}

	p.nextToken(ctx)

	return stm, nil
}

func (p *Parser) parseRowStatement(ctx context.Context) (ast.Statement, error) {
	slogctx.Debug(ctx, "Parser.parseRowStatement()")

	var (
		err  error
		stmt *ast.Assignment
	)

	name := p.token.Literal
	active := !p.token.Commented

	p.nextToken(ctx)

	switch p.token.Type {
	case token.NewLine, token.EOF:
		stmt = p.parseNakedAssign(ctx, name)

	case token.Assign:
		p.nextToken(ctx)

		switch p.token.Type {
		case token.NewLine, token.EOF:
			stmt = p.parseNakedAssign(ctx, name)

		case token.Value, token.RawValue:
			stmt, err = p.parseCompleteAssign(ctx, name)

		default:
			_, err = p.unexpectedToken("parseRowStatement 1")
		}

	default:
		_, err = p.unexpectedToken("parseRowStatement 2")
	}

	if err != nil {
		return nil, err
	}

	if stmt != nil {
		stmt.Enabled = active

		return stmt, err
	}

	return p.unexpectedToken("parseRowStatement 3")
}

func (p *Parser) parseNakedAssign(ctx context.Context, name string) *ast.Assignment {
	slogctx.Debug(ctx, "Parser.parseNakedAssign()")

	defer p.nextToken(ctx)

	return &ast.Assignment{
		Name:    name,
		Enabled: p.token.Commented,
		Quote:   token.NoQuote,
		Position: ast.Position{
			FirstLine: p.token.LineNumber,
			Line:      p.token.LineNumber,
			LastLine:  p.token.LineNumber,
		},
	}
}

func (p *Parser) parseCompleteAssign(ctx context.Context, name string) (*ast.Assignment, error) {
	slogctx.Debug(ctx, "Parser.parseCompleteAssign()")

	assignment := p.token

	p.nextToken(ctx)

	switch p.token.Type {
	case token.NewLine, token.EOF:
		defer p.nextToken(ctx)

		return &ast.Assignment{
			Name:     name,
			Literal:  assignment.Literal,
			Complete: true,
			Enabled:  p.token.Commented,
			Quote:    assignment.Quote,
			Position: ast.Position{
				FirstLine: p.token.LineNumber,
				Line:      p.token.LineNumber,
				LastLine:  p.token.LineNumber,
			},
		}, nil

	default:
		_, err := p.unexpectedToken("parseCompleteAssign 1")

		return nil, err
	}
}

func (p *Parser) nextToken(ctx context.Context) {
	p.previousToken = p.token
	p.token = p.scanner.NextToken(ctx)
}

func (p *Parser) wasEmptyLine() bool {
	if p.token.Type == token.NewLine && p.previousToken.Type == token.NewLine {
		return true
	}

	return p.token.LineNumber != p.previousToken.LineNumber
}

func (p *Parser) skipBlankLine(ctx context.Context) {
	slogctx.Debug(ctx, "Parser.skipBlankLine()")

	for p.token.Type == token.NewLine || p.token.Type == token.Space {
		p.nextToken(ctx)
	}
}

func (p *Parser) unexpectedToken(details string) (ast.Statement, error) {
	return nil, fmt.Errorf("unexpected token at line %d: %s(%s) - %s", p.token.LineNumber, p.token.Type, p.token.Literal, details)
}
