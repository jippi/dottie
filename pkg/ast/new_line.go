package ast

type Newline struct {
	Blank      bool
	LineNumber int
	Group      *Group
}

func (s *Newline) statementNode() {
}
