package token

type Quote uint

const (
	DoubleQuotes Quote = iota
	SingleQuotes
	NoQuotes
)

var quotes = []rune{
	SingleQuotes: '\'',
	DoubleQuotes: '"',
	NoQuotes:     0,
}

func (qt Quote) Is(in rune) bool {
	return quotes[qt] == in
}

func (qt Quote) Rune() rune {
	return quotes[qt]
}

// String returns the string corresponding to the token.
func (qt Quote) String() string {
	s := ""

	if int(qt) < len(quotes) {
		s = string(quotes[qt])
	}

	if s == "" {
		s = "quote(" + string(qt.Rune()) + ")"
	}

	return s
}
