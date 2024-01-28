package ast

import (
	"fmt"
)

// Comment node represents a comment statement.
type Comment struct {
	Value           string
	LineNumber      int
	AnnotationKey   string
	AnnotationValue string
	Group           *Group
}

func NewComment(value string) *Comment {
	return &Comment{
		Value: " " + value,
	}
}

func (s *Comment) statementNode() {
}

func (s *Comment) String() string {
	return fmt.Sprintf("#%s", s.Value)
}
