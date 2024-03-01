//nolint:errname,nlreturn,wsl,varnamelen
package console

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeSplitWords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []Word
	}{
		{
			name:     "empty",
			input:    "",
			expected: nil,
		},
		{
			name:  "hello world",
			input: `hello world`,
			expected: []Word{
				{Start: 0, Stop: 5, Value: "hello"},
				{Start: 6, Stop: 11, Value: "world"},
			},
		},
		{
			name:  "hello world",
			input: `"hello" world`,
			expected: []Word{
				{Start: 0, Stop: 5, Value: "hello"},
				{Start: 6, Stop: 11, Value: "world"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := SafeSplitWords(tt.input)
			require.EqualValues(t, tt.expected, actual)
		})
	}
}
