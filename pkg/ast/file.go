// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"bytes"
)

// Node represents AST-node of the syntax tree.
type Node interface{}

// Statement represents syntax tree node of .env file statement (like: assignment or comment).
type Statement interface {
	Node
	statementNode()
}

// File node represents .env file statement, that contains assignments and comments.
type File struct {
	Statements []Statement
	Groups     []*Group
}

func (s *File) statementNode() {
}

func (s *File) Pairs() map[string]string {
	values := map[string]string{}

	for _, stmt := range s.Statements {
		if assign, ok := stmt.(*Assignment); ok {
			values[assign.Key] = assign.Value
		}
	}

	return values
}

func (s *File) GetGroup(name string) *Group {
	for _, grp := range s.Groups {
		if grp.Name == name {
			return grp
		}
	}

	return nil
}

func (s *File) Get(name string) *Assignment {
	for _, stmt := range s.Statements {
		if assign, ok := stmt.(*Assignment); ok {
			if assign.Key == name {
				return assign
			}
		}
	}

	return nil
}

func (s *File) Render() []byte {
	return s.RenderWithFilter(nil)
}

func (s *File) RenderWithFilter(f *RenderSettings) []byte {
	var buff bytes.Buffer

	for _, s := range s.Statements {
		switch v := s.(type) {
		case *Group:
			if f == nil || f.Groups() {
				buff.WriteString("################################################################################")
				buff.WriteString("\n")

				buff.WriteString("# " + v.Name)
				buff.WriteString("\n")

				buff.WriteString("################################################################################")
				buff.WriteString("\n")
			}

		case *Comment:
			if f == nil || f.Comments() {
				buff.WriteString(v.String())
				buff.WriteString("\n")
			}

		case *Assignment:
			if f != nil && !f.Match(v) {
				continue
			}

			if f == nil || f.Comments() {
				buff.WriteString(v.String())
				buff.WriteString("\n")
				buff.WriteString("\n")

				continue
			}

			buff.WriteString(v.Assignment())
			buff.WriteString("\n")

		case *Newline:
			buff.WriteString("\n")
		}
	}

	return bytes.TrimSpace(buff.Bytes())
}
