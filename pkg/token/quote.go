package token

import (
	"encoding/json"
)

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

var codeMap = map[rune][]rune{
	'\n': []rune(`\n`),
	'\v': []rune(`\v`),
	'\r': []rune(`\r`),
	'\t': []rune(`\t`),
	'\f': []rune(`\f`),
}

func (qt Quote) Escape(value string) string {
	outcome := make([]rune, 0)

	for _, runeVal := range value {
		if runeVal == '\n' || runeVal == '\v' || runeVal == '\r' || runeVal == '\t' || runeVal == '\f' {
			outcome = append(outcome, codeMap[runeVal]...)

			continue
		}

		outcome = append(outcome, runeVal)
	}

	return string(outcome)
}

func (qt Quote) Unescape(value string) string {
	return value

	outcome := make([]rune, 0)

	var prev rune

	for i, runeVal := range value {
		if prev == '\\' && (runeVal == '\n' || runeVal == '\v' || runeVal == '\r' || runeVal == '\t') {
			outcome = outcome[:i-1]
			outcome = append(outcome, []rune(string(runeVal))...)

			continue
		}

		prev = runeVal
		outcome = append(outcome, runeVal)
	}

	// value = strings.ReplaceAll(value, `\n`, "\n")
	// value = strings.ReplaceAll(value, `\r`, "\r")
	// value = strings.ReplaceAll(value, `\t`, "\t")
	// value = strings.ReplaceAll(value, `\v`, "\v")
	// value = strings.ReplaceAll(value, `\`+qt.String(), qt.String())

	return string(outcome)
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
