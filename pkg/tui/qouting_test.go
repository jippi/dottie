package tui_test

import (
	"strconv"
	"testing"

	"github.com/jippi/dottie/pkg/tui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuote(t *testing.T) {
	t.Parallel()

	input := "\n"

	actual := tui.Quote(input)

	assert.Equal(t, "\\n", actual)
}

func TestUnquote(t *testing.T) {
	t.Parallel()

	newlineRune := '\n'

	out, err := tui.Unquote(`\n`, '"', true)
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

			for _, line := range tui.DumpSlice(tt.input) {
				t.Log(line)
			}

			// Ensure expected runes match the runes from the input
			assert.Equal(t, tt.expectedRunes, []rune(tt.input))

			t.Log("-----------------------")
			t.Log("strconv.Quote")
			t.Log("-----------------------")

			qu := strconv.Quote(tt.input)

			for _, line := range tui.DumpSlice(qu[1 : len(qu)-1]) {
				t.Log(line)
			}

			t.Log("-----------------------")
			t.Log("tui.Quote")
			t.Log("-----------------------")

			// Quote the string
			quoted := tui.Quote(tt.input)

			for _, line := range tui.DumpSlice(quoted) {
				t.Log(line)
			}

			// Ensure output matches the expected quoted string
			assert.Equal(t, tt.expectedQuoted, quoted)

			t.Log("-----------------------")
			t.Log("strconv.Unquote")
			t.Log("-----------------------")

			s, err := strconv.Unquote(qu)
			require.NoError(t, err)

			for _, line := range tui.DumpSlice(s) {
				t.Log(line)
			}

			// Unquote the string back
			unquoted, err := tui.Unquote(quoted, '"', true)
			require.NoError(t, err)

			t.Log("-----------------------")
			t.Log("tui.unquoted")
			t.Log("-----------------------")

			for _, line := range tui.DumpSlice(unquoted) {
				t.Log(line)
			}

			// The unquoted string must be equal to the original input
			assert.Equal(t, tt.input, unquoted, "unquoted string does not match original input")

			// Ensure the unquoted string matches the original string at a rune level
			assert.Equal(t, tt.expectedRunes, []rune(unquoted), "unquoted runes does not match expected runes")
		})
	}
}
