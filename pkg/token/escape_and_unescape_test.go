package token_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/jippi/dottie/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEscapeAndUnescape(t *testing.T) {
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

			for _, line := range token.DebugStringSlice(context.TODO(), tt.input) {
				t.Log(line)
			}

			// Ensure expected runes match the runes from the input
			assert.Equal(t, tt.expectedRunes, []rune(tt.input))

			t.Log("-----------------------")
			t.Log("strconv.Quote")
			t.Log("-----------------------")

			strconvQuoted := strconv.Quote(tt.input)

			for _, line := range token.DebugStringSlice(context.TODO(), strconvQuoted[1:len(strconvQuoted)-1]) {
				t.Log(line)
			}

			t.Log("-----------------------")
			t.Log("tui.Quote")
			t.Log("-----------------------")

			// Quote the string
			tuiQuoted := token.Escape(context.TODO(), tt.input)

			for _, line := range token.DebugStringSlice(context.TODO(), tuiQuoted) {
				t.Log(line)
			}

			// Ensure output matches the expected quoted string
			assert.Equal(t, tt.expectedQuoted, tuiQuoted)

			t.Log("-----------------------")
			t.Log("strconv.Unquote")
			t.Log("-----------------------")

			s, err := strconv.Unquote(strconvQuoted)
			require.NoError(t, err)

			for _, line := range token.DebugStringSlice(context.TODO(), s) {
				t.Log(line)
			}

			// Unquote the string back
			unquoted, err := token.Unescape(context.TODO(), tuiQuoted, '"', true)
			require.NoError(t, err)

			t.Log("-----------------------")
			t.Log("tui.unquoted")
			t.Log("-----------------------")

			for _, line := range token.DebugStringSlice(context.TODO(), unquoted) {
				t.Log(line)
			}

			// The unquoted string must be equal to the original input
			assert.Equal(t, tt.input, unquoted, "unquoted string does not match original input")

			// Ensure the unquoted string matches the original string at a rune level
			assert.Equal(t, tt.expectedRunes, []rune(unquoted), "unquoted runes does not match expected runes")
		})
	}
}
