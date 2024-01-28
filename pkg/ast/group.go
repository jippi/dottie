package ast

type Group struct {
	Name       string
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
