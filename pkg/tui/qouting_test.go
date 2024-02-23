package tui_test

import (
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
			input:          "\\",
			expectedQuoted: `\`,
			expectedRunes:  []rune{'\\'},
		},
		{
			name:           "double-slash",
			input:          "\\\\",
			expectedQuoted: `\\`,
			expectedRunes:  []rune{'\\', '\\'},
		},
		{
			name:           "triple-slash",
			input:          "\\\\\\",
			expectedQuoted: `\\\`,
			expectedRunes:  []rune{'\\', '\\', '\\'},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Ensure expected runes match the runes from the input
			assert.Equal(t, tt.expectedRunes, []rune(tt.input))

			// Quote the string
			quoted := tui.Quote(tt.input)

			// Ensure output matches the expected quoted string
			assert.Equal(t, tt.expectedQuoted, quoted)

			// Unquote the string back
			unquoted, err := tui.Unquote(quoted, '"', true)
			require.NoError(t, err)

			// The unquoted string must be equal to the original input
			assert.Equal(t, tt.input, unquoted, "unquoted string does not match original input")

			// Ensure the unquoted string matches the original string at a rune level
			assert.Equal(t, tt.expectedRunes, []rune(unquoted), "unquoted runes does not match expected runes")
		})
	}
}
