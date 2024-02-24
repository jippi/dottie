package token_test

import (
	"context"
	"testing"

	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/stretchr/testify/require"
)

func TestQuote_Unescape(t *testing.T) {
	return

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

			// actual := tt.quote.Unescape(tt.input)
			actual := tt.input

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
