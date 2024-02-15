package token

import "encoding/json"

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

	str := ""

	if int(qt) < len(quotes) {
		str = string(quotes[qt])
	}

	if str == "" {
		str = "quote(" + string(qt.Rune()) + ")"
	}

	return str
}

func (qt Quote) MarshalJSON() ([]byte, error) {
	switch qt {
	case NoQuotes:
		return json.Marshal(nil)

	case SingleQuotes:
		return json.Marshal("single")

	case DoubleQuotes:
		return json.Marshal("double")

	default:
		panic("unknown quote style")
	}
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
