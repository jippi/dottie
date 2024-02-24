package tui_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuote_Unescape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		quote    token.Quote
		input    string
		expected string
	}{
		{
			name:     "flat string",
			quote:    token.DoubleQuotes,
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "tab",
			quote:    token.DoubleQuotes,
			input:    "\\t",
			expected: "\t",
		},
		{
			name:     "newline",
			quote:    token.DoubleQuotes,
			input:    "\\n",
			expected: "\n",
		},
		{
			name:     "many chars",
			quote:    token.DoubleQuotes,
			input:    `my_key="\t"`,
			expected: "my_key=\"\t\"",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual, err := tui.Unquote(context.TODO(), tt.input, '"', true)
			require.NoError(t, err)

			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestQuote_Escape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		quote    token.Quote
		input    string
		expected string
	}{
		{
			name:     "flat string",
			quote:    token.DoubleQuotes,
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "tab",
			quote:    token.DoubleQuotes,
			input:    "\t",
			expected: "\\t",
		},
		{
			name:     "newline",
			quote:    token.DoubleQuotes,
			input:    "\n",
			expected: "\\n",
		},
		{
			name:     "many chars",
			quote:    token.DoubleQuotes,
			input:    "'           '",
			expected: `'           '`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := tui.Quote(context.TODO(), tt.input)

			require.EqualValues(t, tt.expected, actual)
		})
	}
}

func TestQuote(t *testing.T) {
	t.Parallel()

	input := "\n"

	actual := tui.Quote(context.TODO(), input)

	assert.Equal(t, "\\n", actual)
}

func TestUnquote(t *testing.T) {
	t.Parallel()

	newlineRune := '\n'

	out, err := tui.Unquote(context.TODO(), `\n`, '"', true)
	require.NoError(t, err)
	assert.Equal(t, []rune{newlineRune}, []rune(out))
}

func TestBackAndForth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          string
		expectedQuoted string
		expectedRunes  []rune
	}{
		{
			name:           "newline",
			input:          "\n",
			expectedQuoted: `\n`,
			expectedRunes:  []rune{'\n'},
		},
		{
			name:           "slash",
			input:          `\`,
			expectedQuoted: `\\`,
			expectedRunes:  []rune{'\\'},
		},
		{
			name:           "double-slash",
			input:          `\\`,
			expectedQuoted: `\\\\`,
			expectedRunes:  []rune{'\\', '\\'},
		},
		{
			name:           "triple-slash",
			input:          `\\\`,
			expectedQuoted: `\\\\\\`,
			expectedRunes:  []rune{'\\', '\\', '\\'},
		},
		{
			name:           "null",
			input:          "\x00",
			expectedQuoted: `\x00`,
			expectedRunes:  []rune{0},
		},
		{
			name:           "slash-zero",
			input:          "\\0",
			expectedQuoted: `\\0`,
			expectedRunes:  []rune{'\\', '0'},
		},
		{
			name:           "weird-1",
			input:          "\xf5",
			expectedQuoted: "\\xf5",
			expectedRunes:  []rune{65533},
		},
		{
			name:           "weird-2",
			input:          "\x00$",
			expectedQuoted: "\\x00$",
			expectedRunes:  []rune{0, 36},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Log("-----------------------")
			t.Log("tt.input")
			t.Log("-----------------------")

			for _, line := range tui.DumpSlice(context.TODO(), tt.input) {
				t.Log(line)
			}

			// Ensure expected runes match the runes from the input
			assert.Equal(t, tt.expectedRunes, []rune(tt.input))

			t.Log("-----------------------")
			t.Log("strconv.Quote")
			t.Log("-----------------------")

			strconvQuoted := strconv.Quote(tt.input)

			for _, line := range tui.DumpSlice(context.TODO(), strconvQuoted[1:len(strconvQuoted)-1]) {
				t.Log(line)
			}

			t.Log("-----------------------")
			t.Log("tui.Quote")
			t.Log("-----------------------")

			// Quote the string
			tuiQuoted := tui.Quote(context.TODO(), tt.input)

			for _, line := range tui.DumpSlice(context.TODO(), tuiQuoted) {
				t.Log(line)
			}

			// Ensure output matches the expected quoted string
			assert.Equal(t, tt.expectedQuoted, tuiQuoted)

			t.Log("-----------------------")
			t.Log("strconv.Unquote")
			t.Log("-----------------------")

			s, err := strconv.Unquote(strconvQuoted)
			require.NoError(t, err)

			for _, line := range tui.DumpSlice(context.TODO(), s) {
				t.Log(line)
			}

			// Unquote the string back
			unquoted, err := tui.Unquote(context.TODO(), tuiQuoted, '"', true)
			require.NoError(t, err)

			t.Log("-----------------------")
			t.Log("tui.unquoted")
			t.Log("-----------------------")

			for _, line := range tui.DumpSlice(context.TODO(), unquoted) {
				t.Log(line)
			}

			// The unquoted string must be equal to the original input
			assert.Equal(t, tt.input, unquoted, "unquoted string does not match original input")

			// Ensure the unquoted string matches the original string at a rune level
			assert.Equal(t, tt.expectedRunes, []rune(unquoted), "unquoted runes does not match expected runes")
		})
	}
}
