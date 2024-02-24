package token_test

import (
	"context"
	"testing"

	"github.com/jippi/dottie/pkg/token"
	"github.com/stretchr/testify/require"
)

func TestUnescape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		quote    token.Quote
		input    string
		expected string
	}{
		{
			name:     "flat string",
			quote:    token.DoubleQuote,
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "tab",
			quote:    token.DoubleQuote,
			input:    "\\t",
			expected: "\t",
		},
		{
			name:     "newline",
			quote:    token.DoubleQuote,
			input:    `\n`,
			expected: "\n",
		},
		{
			name:     "many chars",
			quote:    token.DoubleQuote,
			input:    `my_key="\t"`,
			expected: "my_key=\"\t\"",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual, err := token.Unescape(context.TODO(), tt.input, token.DoubleQuote)
			require.NoError(t, err)

			require.Equal(t, tt.expected, actual)
		})
	}
}
