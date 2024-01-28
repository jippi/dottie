package ast

type Group struct {
	Comment    string
	FirstLine  int
	LastLine   int
	Statements []Statement
}

func (s *Group) statementNode() {
}

type Newline struct {
	Blank      bool
	LineNumber int
	Group      *Group
}

func (s *Newline) statementNode() {
}
