package token

import (
	"encoding/json"
	"errors"
)

var ErrInvalidQuoteStyle = errors.New("Invalid quote style")

type Quote uint

const (
	InvalidQuote Quote = iota
	DoubleQuote
	SingleQuote
	NoQuote

	maxQuote
)

var quotes = []rune{
	SingleQuote: '\'',
	DoubleQuote: '"',
	NoQuote:     0,
}

func (qt Quote) Is(in rune) bool {
	return quotes[qt] == in
}

func (qt Quote) Valid() bool {
	return qt > InvalidQuote && qt < maxQuote
}

func (qt Quote) Rune() rune {
	return quotes[qt]
}

func (qt Quote) Byte() byte {
	return byte(quotes[qt])
}

// String returns the string corresponding to the token.
func (qt Quote) String() string {
	// the NoQuotes rune (0) are *not* the same as an empty string, so we handle it specially here
	if qt == NoQuote {
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
	case NoQuote:
		return json.Marshal(nil)

	case SingleQuote:
		return json.Marshal("single")

	case DoubleQuote:
		return json.Marshal("double")

	default:
		panic("unknown quote style")
	}
}

func QuoteFromString(in string) Quote {
	switch in {
	case "\"", "double":
		return DoubleQuote

	case "'", "single":
		return SingleQuote

	case "none":
		return NoQuote

	default:
		return InvalidQuote
	}
}
