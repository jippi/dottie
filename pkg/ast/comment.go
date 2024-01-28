package ast

// Comment node represents a comment statement.
type Comment struct {
	Value           string
	LineNumber      int
	Annotation      bool
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
	return s.Value
}
