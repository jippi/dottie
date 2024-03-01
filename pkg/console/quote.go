//nolint:errname,nlreturn,wsl,varnamelen
package console

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	UnterminatedSingleQuoteError = errors.New("Unterminated single-quoted string")
	UnterminatedDoubleQuoteError = errors.New("Unterminated double-quoted string")
	UnterminatedEscapeError      = errors.New("Unterminated backslash-escape")
)

var (
	splitChars        = " \n\t"
	singleChar        = '\''
	doubleChar        = '"'
	escapeChar        = '\\'
	doubleEscapeChars = "$`\"\n\\"
)

type Word struct {
	Start int
	Stop  int
	Value string
}

func (w Word) String() string {
	return fmt.Sprintf("%s (%d:%d) | ", w.Value, w.Start, w.Stop)
}

func JoinWords(in []Word) []string {
	out := make([]string, len(in))

	for i, w := range in {
		out[i] = w.Value
	}

	return out
}

func SafeSplitWords(input string) []Word {
	words, err := Split(input)
	if errors.Is(err, UnterminatedDoubleQuoteError) {
		return SafeSplitWords(input + `"`)
	}

	if errors.Is(err, UnterminatedSingleQuoteError) {
		return SafeSplitWords(input + `'`)
	}

	if err != nil {
		panic(err)
	}

	return words
}

// Split splits a string according to /bin/sh's word-splitting rules. It
// supports backslash-escapes, single-quotes, and double-quotes. Notably it does
// not support the $” style of quoting. It also doesn't attempt to perform any
// other sort of expansion, including brace expansion, shell expansion, or
// pathname expansion.
//
// If the given input has an unterminated quoted string or ends in a
// backslash-escape, one of UnterminatedSingleQuoteError,
// UnterminatedDoubleQuoteError, or UnterminatedEscapeError is returned.
func Split(input string) (words []Word, err error) {
	var (
		offset int
		tmp    int
	)

	for len(input) > 0 {
		// skip any splitChars at the start
		c, l := utf8.DecodeRuneInString(input)

		if strings.ContainsRune(splitChars, c) {
			input = input[l:]

			continue
		}

		if c == escapeChar {
			// Look ahead for escaped newline so we can skip over it
			next := input[l:]
			if len(next) == 0 {
				err = UnterminatedEscapeError

				return
			}

			c2, l2 := utf8.DecodeRuneInString(next)
			if c2 == '\n' {
				input = next[l2:]

				continue
			}
		}

		var word Word

		word, input, err = splitWord(input)
		if err != nil {
			return
		}

		tmp = offset + word.Stop

		word.Start = word.Start + offset
		word.Stop = word.Stop + offset

		offset = tmp + 1

		words = append(words, word)
	}

	return
}

func splitWord(input string) (word Word, remainder string, err error) {
	var (
		buf   bytes.Buffer
		start = 0
		stop  = len(input)
	)

raw:
	{
		cur := input

		for len(cur) > 0 {
			c, l := utf8.DecodeRuneInString(cur)
			cur = cur[l:]

			if c == singleChar {
				start = 0
				stop = len(input) - len(cur) - l

				buf.WriteString(input[start:stop])
				input = cur

				goto single
			}

			if c == doubleChar {
				start = 0
				stop = len(input) - len(cur) - l

				buf.WriteString(input[start:stop])
				input = cur

				goto double
			}

			if c == escapeChar {
				start = 0
				stop = len(input) - len(cur) - l

				buf.WriteString(input[start:stop])
				input = cur

				goto escape
			}

			if strings.ContainsRune(splitChars, c) {
				start := 0
				stop := len(input) - len(cur) - l

				buf.WriteString(input[start:stop])

				goto done
			}
		}

		if len(input) > 0 {
			word.Start = 0
			word.Stop = len(input)

			buf.WriteString(input)
		}

		input = ""
		fmt.Println("done?!")
		goto done
	}

escape:
	{
		if len(input) == 0 {
			return Word{}, "", UnterminatedEscapeError
		}

		c, l := utf8.DecodeRuneInString(input)
		if c == '\n' {
			// a backslash-escaped newline is elided from the output entirely
		} else {
			start = 0
			stop = l

			buf.WriteString(input[:l])
		}

		input = input[l:]
	}

	goto raw

single:
	{
		i := strings.IndexRune(input, singleChar)
		if i == -1 {
			return word, "", UnterminatedSingleQuoteError
		}

		start = 0
		stop = i

		buf.WriteString(input[start:stop])
		input = input[i+1:]

		goto raw
	}

double:
	{
		cur := input

		for len(cur) > 0 {
			c, l := utf8.DecodeRuneInString(cur)
			cur = cur[l:]
			if c == doubleChar {
				start = 0
				stop = len(input) - len(cur) - l

				buf.WriteString(input[start:stop])
				input = cur

				goto raw
			}

			if c == escapeChar {
				// bash only supports certain escapes in double-quoted strings
				c2, l2 := utf8.DecodeRuneInString(cur)

				cur = cur[l2:]

				if strings.ContainsRune(doubleEscapeChars, c2) {
					start = 0
					stop = len(input) - len(cur) - l - l2

					buf.WriteString(input[start:stop])

					if c2 == '\n' {
						// newline is special, skip the backslash entirely
					} else {
						buf.WriteRune(c2)
					}

					input = cur
				}
			}
		}

		word.Start = 0
		word.Stop = len(input)

		return word, "", UnterminatedDoubleQuoteError
	}

done:

	return Word{
		Value: buf.String(),
		Start: start,
		Stop:  stop,
	}, input, nil
}
