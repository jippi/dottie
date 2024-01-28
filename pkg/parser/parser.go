// Package parser implements a parser for the .env files.
package parser

import (
	"fmt"
	"strings"

	"dotfedi/pkg/ast"
	"dotfedi/pkg/token"
)

// Scanner converts a sequence of characters into a sequence of tokens.
type Scanner interface {
	NextToken() token.Token
}

// Parser takes a Scanner and builds an abstract syntax tree.
type Parser struct {
	scanner       Scanner
	token         token.Token
	previousToken token.Token
}

// New returns new Parser.
func New(scanner Scanner) *Parser {
	return &Parser{
		scanner: scanner,
		token:   scanner.NextToken(),
	}
}

// Parse parses the .env file and returns an ast.Statement.
func (p *Parser) Parse() (ast.Statement, error) {
	// var statements []ast.Statement
	var currentGroup *ast.Group
	var comments []*ast.Comment
	var previousStatement ast.Statement

	result := &ast.File{}

	for p.token.Type != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		switch val := stmt.(type) {
		case *ast.Group:
			// Track the last line of this group
			if currentGroup != nil {
				currentGroup.LastLine = p.token.LineNumber
			}

			// Change the group
			currentGroup = val

			// Append the group
			result.Groups = append(result.Groups, currentGroup)

			// Append it to the statements list
			result.Statements = append(result.Statements, val)

		case *ast.Assignment:
			// Assign the assignment to a grouping if such exists
			if currentGroup != nil {
				currentGroup.Statements = append(currentGroup.Statements, val)

				val.Group = currentGroup
			}

			// Assign accumulated comments to this assignment
			val.Comments = comments

			if len(val.Comments) > 0 {
				val.FirstLine = val.Comments[0].LineNumber
			}

			val.LastLine = val.LineNumber

			// Reset comment block
			comments = nil

			result.Statements = append(result.Statements, stmt)

		case *ast.Comment:
			if currentGroup != nil {
				val.Group = currentGroup
			}

			comments = append(comments, val)

		case *ast.Newline:
			switch previousStatement.(type) {
			// If the previous statement was an assignment, ignore the newline
			// as we will be emitted that ourself later
			case *ast.Assignment:
				continue
			}

			if !val.Blank {
				continue
			}

			// If there is a blank line, print all previous comments
			for _, c := range comments {
				result.Statements = append(result.Statements, c)
			}

			// Reset the accumulated comments slice
			comments = nil

			// Attach the newline to a group for easier filtering
			if currentGroup != nil {
				val.Group = currentGroup
			}

			result.Statements = append(result.Statements, val)
		}

		previousStatement = stmt
	}

	if currentGroup != nil {
		currentGroup.LastLine = p.token.LineNumber
	}

	return result, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.token.Type {
	case token.Identifier:
		return p.parseRowStatement()

	case token.Comment, token.CommentAnnotation:
		return p.parseCommentStatement()

	case token.EOF:
		return nil, nil

	case token.NewLine:
		res := &ast.Newline{
			Blank:      p.wasEmptyLine(),
			LineNumber: p.token.LineNumber,
		}

		p.nextToken()

		return res, nil

	default:
		panic(fmt.Errorf("(B) unexpected statement: %s(%q)", p.token.Type, p.token.Literal))
	}
}

func (p *Parser) parseCommentStatement() (ast.Statement, error) {
	// If the comment doesn't look like a header, just treat it as a normal comment and move on
	if !strings.Contains(p.token.Literal, "###") {
		stm := &ast.Comment{
			Value:           p.token.Literal,
			LineNumber:      p.token.LineNumber,
			Annotation:      p.token.Annotation,
			AnnotationKey:   p.token.AnnotationKey,
			AnnotationValue: p.token.AnnotationValue,
		}

		p.nextToken()

		return stm, nil
	}

	// If the comment block looks like a header group
	group := &ast.Group{
		FirstLine: p.token.LineNumber,
	}

	// Move forward
	p.nextToken()
	p.skipBlankLine()

	switch p.token.Type {
	case token.Comment:
		break

	default:
		panic("invalid")
	}
	group.Name = strings.TrimSpace(p.token.Literal)

	p.skipGroupHeader()

	return group, nil
}

func (p *Parser) parseRowStatement() (ast.Statement, error) {
	var err error
	var stmt *ast.Assignment

	name := p.token.Literal
	shadow := p.token.Commented

	p.nextToken()

	switch p.token.Type {
	case token.NewLine, token.EOF:
		stmt, err = p.parseNakedAssign(name)

	case token.Assign:
		p.nextToken()

		switch p.token.Type {
		case token.NewLine, token.EOF:
			stmt, err = p.parseNakedAssign(name)

		case token.Value, token.RawValue:
			stmt, err = p.parseCompleteAssign(name)
		}
	}

	if err != nil {
		return nil, err
	}

	if stmt != nil {
		stmt.Commented = shadow

		return stmt, err
	}

	return nil, fmt.Errorf("unexpected token at line %d: %s(%s)", p.token.LineNumber, p.token.Type, p.token.Literal)
}

func (p *Parser) parseNakedAssign(name string) (*ast.Assignment, error) {
	defer p.nextToken()

	return &ast.Assignment{
		Key:        name,
		LineNumber: p.token.LineNumber,
		Naked:      true,
		Commented:  p.token.Commented,
	}, nil
}

func (p *Parser) parseCompleteAssign(name string) (*ast.Assignment, error) {
	value := p.token.Literal
	quoted := p.token.Quoted
	p.nextToken()

	switch p.token.Type {
	case token.NewLine, token.EOF:
		p.nextToken()

		return &ast.Assignment{
			Key:        name,
			Value:      value,
			LineNumber: p.token.LineNumber,
			Complete:   true,
			Commented:  p.token.Commented,
			Quoted:     quoted,
		}, nil

	default:
		return nil, fmt.Errorf("unexpected token at line %d: %s(%s)", p.token.LineNumber, p.token.Type, p.token.Literal)
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

func (p *Parser) skipGroupHeader() {
	p.nextToken()

	p.skipBlankLine()

	switch p.token.Type {
	case token.Comment:
		break

	default:
		panic("invalid")
	}

	p.nextToken()
}
