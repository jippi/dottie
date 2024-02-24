package tui

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	slogctx "github.com/veqryn/slog-context"
)

func Unquote(ctx context.Context, input string, quote byte, unescape bool) (out string, err error) {
	input0 := input

	// Handle quoted strings without any escape sequences.
	if !contains(input, '\\') && !contains(input, '\n') {
		var valid bool

		switch quote {
		case '"':
			valid = utf8.ValidString(input)

		case '\'':
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
		slogctx.Debug(ctx, "Unquote.input", StringDump(input))

		// Process the next character, rejecting any unescaped newline characters which are invalid.
		runeVal, multibyte, remaining, err := UnquoteChar(ctx, input, quote)
		if err != nil {
			return input0, err
		}

		slogctx.Debug(ctx, "Unquote.remaining", StringDump(remaining))

		input = remaining

		// Append the character if unescaping the input.
		if unescape {
			if runeVal < utf8.RuneSelf || !multibyte {
				slogctx.Debug(ctx, "-- append")

				buf = append(buf, byte(runeVal))
			} else {
				slogctx.Debug(ctx, "-- utf8.AppendRune")

				buf = utf8.AppendRune(buf, runeVal)
			}
		}

		// Single quoted strings must be a single character.
		if quote == '\'' {
			break
		}
	}

	if unescape {
		return string(buf), nil
	}

	return input, nil
}

func contains(s string, c byte) bool {
	return index(s, c) != -1
}

func index(s string, c byte) int {
	return strings.IndexByte(s, c)
}

func UnquoteChar(ctx context.Context, input string, quote byte) (value rune, multibyte bool, tail string, err error) {
	slogctx.Debug(ctx, "UnquoteChar.start.input", StringDump(input))

	if len(input) == 0 {
		return
	}

	slogctx.Debug(ctx, "UnquoteChar.switch.char", StringDump(string(input[0])))

	switch char := input[0]; {
	case char >= utf8.RuneSelf:
		slogctx.Debug(ctx, "UnquoteChar.switch.outcome: char >= utf8.RuneSelf")

		r, size := utf8.DecodeRuneInString(input)

		return r, true, input[size:], nil

	case char != '\\':
		slogctx.Debug(ctx, "UnquoteChar.switch.outcome: char != '\\'")

		return rune(input[0]), false, input[1:], nil

	case char == '\\' && len(input) <= 1:
		slogctx.Debug(ctx, "UnquoteChar.switch.len <= 1", StringDump(string(input[0])), StringDump(input[1:]))

		return rune(input[0]), false, input[1:], nil
	}

	slogctx.Debug(ctx, "UnquoteChar.switch.miss: yup")

	// initial := input[0]

	char := input[1]
	input = input[2:]

	slogctx.Debug(ctx, "UnquoteChar.char", StringDump(string(char)))
	slogctx.Debug(ctx, "UnquoteChar.input", StringDump(input))

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
		slogctx.Debug(ctx, "UnquoteChar.char-switch: xuU")

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
			slogctx.Debug(ctx, "UnquoteChar: len(s) < n")

			return
		}

		for j := 0; j < n; j++ {
			x, ok := unhex(input[j])
			if !ok {
				err = errors.New("UnquoteChar: unhex error")

				return
			}

			v = v<<4 | x
		}

		input = input[n:]

		if char == 'x' {
			slogctx.Debug(ctx, "UnquoteChar.char-switch: xuU -> x (NOT UTF-8)")

			// single-byte string, possibly not UTF-8
			value = v

			break
		}

		if !utf8.ValidRune(v) {
			err = errors.New("UnquoteChar: invalid rune")

			return
		}

		value = v
		multibyte = true

	case '0', '1', '2', '3', '4', '5', '6', '7':
		slogctx.Debug(ctx, "UnquoteChar.char-switch: numbers")

		v := rune(char) - '0'

		if len(input) < 2 {
			value = v

			// err = errors.New("UnquoteChar: len(s) < 2")

			return
		}

		for j := 0; j < 2; j++ { // one digit already; two more
			x := rune(input[j]) - '0'
			if x < 0 || x > 7 {
				err = errors.New("UnquoteChar: x < 0 || x > 7")

				return
			}

			v = (v << 3) | x
		}

		input = input[2:]

		if v > 255 {
			err = errors.New("UnquoteChar: v > 255")

			return
		}

		value = v

	case '\\':
		value = '\\'

	case '\'':
		if char != quote {
			value = rune(char)
		} else {
			err = errors.New("UnquoteChar single: c != quote")
		}

	case '"':
		if char != quote {
			err = errors.New("UnquoteChar double: c != quote")

			return
		}

		value = rune(char)

	default:
		err = errors.New("UnquoteChar: default: >" + fmt.Sprintf("%U", []rune(string(char))) + "< aka >" + fmt.Sprintf("%q", char) + "<")

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
