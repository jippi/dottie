// nolint varnamelen
package token

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/jippi/dottie/pkg/tui"
	slogctx "github.com/veqryn/slog-context"
)

func Unescape(ctx context.Context, input string, quote Quote) (out string, err error) {
	if !quote.Valid() {
		panic(ErrInvalidQuoteStyle)
	}

	ctx = slogctx.With(
		ctx,
		slog.String("source", "token.Unescape()"),
		slog.String("quote", quote.Name()),
	)

	input0 := input

	// Handle quoted strings without any escape sequences.
	if !contains(input, '\\') && !contains(input, '\n') {
		var valid bool

		switch quote {
		case DoubleQuote:
			valid = utf8.ValidString(input)

		case SingleQuote:
			r, n := utf8.DecodeRuneInString(input)
			valid = (r != utf8.RuneError || n != 1)
		}

		if valid {
			return input, nil
		}
	}

	var buf []byte

	// LOOP
	for len(input) > 0 {
		slogctx.Debug(ctx, "Unescape :: loop :: input", tui.StringDump("input", input))

		// Process the next character, rejecting any unescaped newline characters which are invalid.
		runeVal, multibyte, remaining, err := unescapeChar(ctx, input, quote)
		if err != nil {
			return input0, err
		}

		slogctx.Debug(ctx, "Unescape :: loop :: remaining", tui.StringDump("remaining", remaining))

		input = remaining

		if runeVal < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(runeVal))
		} else {
			buf = utf8.AppendRune(buf, runeVal)
		}

		// Single quoted strings must be a single character.
		if quote == SingleQuote {
			break
		}
	}

	return string(buf), nil
}

func contains(s string, c byte) bool {
	return index(s, c) != -1
}

func index(s string, c byte) int {
	return strings.IndexByte(s, c)
}

func unescapeChar(ctx context.Context, input string, quote Quote) (value rune, multibyte bool, tail string, err error) {
	ctx = slogctx.With(ctx, tui.StringDump("input", input))
	slogctx.Debug(ctx, "token.unescapeChar()")

	if len(input) == 0 {
		return
	}

	switch char := input[0]; {
	case char >= utf8.RuneSelf:
		r, size := utf8.DecodeRuneInString(input)

		return r, true, input[size:], nil

	case char != '\\':
		return rune(input[0]), false, input[1:], nil

	case char == '\\' && len(input) <= 1:
		return rune(input[0]), false, input[1:], nil
	}

	char := input[1]
	input = input[2:]

	ctx = slogctx.With(
		ctx,
		tui.StringDump("char", string(char)),
		tui.StringDump("input", input),
	)

	slogctx.Debug(ctx, "token.unescapeChar() complex unescape path")

	switch char {
	case 'a':
		value = '\a'

	case 'b':
		value = '\b'

	case 'f':
		value = '\f'

	case 'n':
		value = '\n'

	case 'r':
		value = '\r'

	case 't':
		value = '\t'

	case 'v':
		value = '\v'

	case 'x', 'u', 'U':
		n := 0

		switch char {
		case 'x':
			n = 2

		case 'u':
			n = 4

		case 'U':
			n = 8
		}

		var v rune

		if len(input) < n {
			slogctx.Debug(ctx, "UnescapeChar: len(s) < n")

			return
		}

		for j := 0; j < n; j++ {
			x, ok := unhex(input[j])
			if !ok {
				err = errors.New("UnescapeChar: unhex error")

				return
			}

			v = v<<4 | x
		}

		input = input[n:]

		if char == 'x' {
			slogctx.Debug(ctx, "UnescapeChar.char-switch: xuU -> x (NOT UTF-8)")

			// single-byte string, possibly not UTF-8
			value = v

			break
		}

		if !utf8.ValidRune(v) {
			err = errors.New("UnescapeChar: invalid rune")

			return
		}

		value = v
		multibyte = true

	case '0', '1', '2', '3', '4', '5', '6', '7':
		v := rune(char) - '0'

		if len(input) < 2 {
			value = v

			// err = errors.New("UnescapeChar: len(s) < 2")

			return
		}

		for j := 0; j < 2; j++ { // one digit already; two more
			x := rune(input[j]) - '0'
			if x < 0 || x > 7 {
				err = errors.New("UnescapeChar: x < 0 || x > 7")

				return
			}

			v = (v << 3) | x
		}

		input = input[2:]

		if v > 255 {
			err = errors.New("UnescapeChar: v > 255")

			return
		}

		value = v

	case '\\':
		value = '\\'

	case '\'':
		if char != quote.Byte() {
			value = rune(char)
		} else {
			err = errors.New("UnescapeChar single: c != quote")
		}

	case '"':
		if char != quote.Byte() {
			err = errors.New("UnescapeChar double: c != quote")

			return
		}

		value = rune(char)

	default:
		err = errors.New("UnescapeChar: default: >" + fmt.Sprintf("%U", []rune(string(char))) + "< aka >" + fmt.Sprintf("%q", char) + "<")

		return
	}

	tail = input

	return
}

func unhex(b byte) (v rune, ok bool) {
	c := rune(b)

	switch {
	case '0' <= c && c <= '9':
		return c - '0', true

	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true

	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return
}
