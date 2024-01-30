package token

type Token struct {
	Type       Type
	Literal    string
	Offset     int
	Length     int
	LineNumber uint
	Commented  bool
	Quote      Quote
	Annotation *Annotation
}

func New(t Type, options ...Option) Token {
	token := &Token{
		Type:    t,
		Literal: t.String(),
	}

	for _, o := range options {
		o(token)
	}

	token.Offset = token.Offset - token.Length

	return *token
}
