package parser

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/token"
)

// Scanner converts a sequence of characters into a sequence of tokens.
type Scanner interface {
	NextToken() token.Token
}

// Parser takes a Scanner and builds an abstract syntax tree.
type Parser struct {
	filename      string
	scanner       Scanner
	token         token.Token
	previousToken token.Token
}

// New returns new Parser.
func New(scanner Scanner, filename string) *Parser {
	return &Parser{
		filename: filename,
		scanner:  scanner,
		token:    scanner.NextToken(),
	}
}

// Parse parses the .env file and returns an ast.Statement.
func (p *Parser) Parse(_ context.Context) (document *ast.Document, err error) {
	defer func() {
		if recoveryErr := recover(); recoveryErr != nil {
			err = fmt.Errorf("panic from parser: %+v", recoveryErr)
			document = nil
		}
	}()

	var (
		comments          []*ast.Comment
		currentGroup      *ast.Group
		previousStatement ast.Statement
		statementIndex    int
	)

	document = ast.NewDocument()

	for p.token.Type != token.EOF {
		fmt.Println("WORKING ON THIS TOKEN", spew.Sdump(p.token))

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		fmt.Println("STMT", spew.Sdump(stmt))

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
			comments = nil
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
			comments = nil
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
	document.Initialize()

	return document, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.token.Type {
	case token.Identifier:
		return p.parseRowStatement()

	case token.Comment, token.CommentAnnotation:
		return p.parseCommentStatement()

	case token.GroupBanner:
		return p.parseGroupStatement()

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

		p.nextToken()

		return res, nil

	default:
		return nil, fmt.Errorf("(B) unexpected statement: %s(%q)", p.token.Type, p.token.Literal)
	}
}

func (p *Parser) parseGroupStatement() (ast.Statement, error) {
	var group *ast.Group

	p.nextToken()
	p.skipBlankLine()

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

	p.nextToken()
	p.skipBlankLine()

	switch p.token.Type {
	case token.GroupBanner:
		p.nextToken()

		return group, nil

	default:
		return p.unexpectedToken("parseGroupStatement 2")
	}
}

func (p *Parser) parseCommentStatement() (ast.Statement, error) {
	stm := &ast.Comment{
		Value:      p.token.Literal,
		Annotation: p.token.Annotation,
		Position: ast.Position{
			Line:      p.token.LineNumber,
			FirstLine: p.token.LineNumber,
			LastLine:  p.token.LineNumber,
		},
	}

	p.nextToken()

	return stm, nil
}

func (p *Parser) parseRowStatement() (ast.Statement, error) {
	var (
		err  error
		stmt *ast.Assignment
	)

	name := p.token.Literal
	active := !p.token.Commented

	p.nextToken()

	fmt.Println("parseRowStatement, next token", spew.Sdump(p.token))

	switch p.token.Type {
	case token.NewLine, token.EOF:
		stmt = p.parseNakedAssign(name)

	case token.Assign:
		p.nextToken()

		fmt.Println("parseRowStatement, after assign", spew.Sdump(p.token))

		switch p.token.Type {
		case token.NewLine, token.EOF:
			stmt = p.parseNakedAssign(name)

		case token.Value, token.RawValue:
			stmt, err = p.parseCompleteAssign(name)

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

func (p *Parser) parseNakedAssign(name string) *ast.Assignment {
	defer p.nextToken()

	return &ast.Assignment{
		Name:    name,
		Enabled: p.token.Commented,
		Quote:   token.NoQuotes,
		Position: ast.Position{
			FirstLine: p.token.LineNumber,
			Line:      p.token.LineNumber,
			LastLine:  p.token.LineNumber,
		},
	}
}

func (p *Parser) parseCompleteAssign(name string) (*ast.Assignment, error) {
	assignment := p.token

	p.nextToken()

	switch p.token.Type {
	case token.NewLine, token.EOF:
		defer p.nextToken()

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

func (p *Parser) nextToken() {
	p.previousToken = p.token
	p.token = p.scanner.NextToken()
}

func (p *Parser) wasEmptyLine() bool {
	if p.token.Type == token.NewLine && p.previousToken.Type == token.NewLine {
		return true
	}

	return p.token.LineNumber != p.previousToken.LineNumber
}

func (p *Parser) skipBlankLine() {
	for p.token.Type == token.NewLine || p.token.Type == token.Space {
		p.nextToken()
	}
}

func (p *Parser) unexpectedToken(details string) (ast.Statement, error) {
	return nil, fmt.Errorf("unexpected token at line %d: %s(%s) - %s", p.token.LineNumber, p.token.Type, p.token.Literal, details)
}
