package token

type Option func(*Token)

func WithLiteral(in string) Option {
	return func(t *Token) {
		length := len(in)

		t.Literal = in
		t.Length = length
	}
}

func WithLiteralRune(in rune) Option {
	str := string(in)

	return func(t *Token) {
		length := len(str)

		t.Literal = str
		t.Length = length
	}
}

func WithLineNumber(in uint) Option {
	return func(t *Token) {
		t.LineNumber = in
	}
}

func WithQuoteType(in QuoteType) Option {
	return func(t *Token) {
		t.QuoteType = in
	}
}

func WithOffset(in int) Option {
	return func(t *Token) {
		t.Offset = in
	}
}

func WithAnnotation(key, value string) Option {
	return func(t *Token) {
		t.Annotation = &Annotation{
			Key:   key,
			Value: value,
		}
	}
}
