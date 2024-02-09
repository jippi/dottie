package token

type Quote uint

const (
	InvalidQuotes Quote = iota
	DoubleQuotes
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

func (qt Quote) Valid() bool {
	return qt > 0
}

func (qt Quote) Rune() rune {
	return quotes[qt]
}

// String returns the string corresponding to the token.
func (qt Quote) String() string {
	// the NoQuotes rune (0) are *not* the same as an empty string, so we handle it specially here
	if qt == NoQuotes {
		return ""
	}

	s := ""

	if int(qt) < len(quotes) {
		s = string(quotes[qt])
	}

	if s == "" {
		s = "quote(" + string(qt.Rune()) + ")"
	}

	return s
}

func QuoteFromString(in string) Quote {
	switch in {
	case "\"", "double":
		return DoubleQuotes

	case "'", "single":
		return SingleQuotes

	case "none":
		return NoQuotes

	default:
		return InvalidQuotes
	}
}
