package token_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEscape(t *testing.T) {
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
			input:    "\t",
			expected: "\\t",
		},
		{
			name:     "newline",
			quote:    token.DoubleQuote,
			input:    "\n",
			expected: "\\n",
		},
		{
			name:     "many chars",
			quote:    token.DoubleQuote,
			input:    "           ",
			expected: `           `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := token.Escape(t.Context(), tt.input, token.DoubleQuote)

			require.EqualValues(t, tt.expected, actual)
		})
	}
}

type quoteTest struct {
	in      string
	out     string
	ascii   string
	graphic string
}

var quotetests = []quoteTest{
	{"\a\b\f\r\n\t\v", `\a\b\f\r\n\t\v`, `\a\b\f\r\n\t\v`, `\a\b\f\r\n\t\v`},
	{"\\", `\\`, `\\`, `\\`},
	{"abc\xffdef", `abc\xffdef`, `abc\xffdef`, `abc\xffdef`},
	{"\u263a", `☺`, `\u263a`, `☺`},
	{"\U0010ffff", `\U0010ffff`, `\U0010ffff`, `\U0010ffff`},
	{"\x04", `\x04`, `\x04`, `\x04`},
	// Some non-printable but graphic runes. Final column is double-quoted.
	{"!\u00a0!\u2000!\u3000!", `!\u00a0!\u2000!\u3000!`, `!\u00a0!\u2000!\u3000!`, "!\u00a0!\u2000!\u3000!"},
	{"\x7f", `\x7f`, `\x7f`, `\x7f`},
}

func TestEscapeFromGo(t *testing.T) {
	t.Parallel()

	for _, tt := range quotetests {
		out := token.Escape(t.Context(), tt.in, token.DoubleQuote)
		assert.Equal(t, tt.out, out)
	}
}

func TestEscapeFromGoASCII(t *testing.T) {
	t.Parallel()

	for _, tt := range quotetests {
		out := token.EscapeFull(t.Context(), tt.in, token.DoubleQuote, true, false)
		assert.Equal(t, tt.ascii, out)
	}
}

func TestEscapeFromGoGraphic(t *testing.T) {
	t.Parallel()

	for _, tt := range quotetests {
		out := token.EscapeFull(t.Context(), tt.in, token.DoubleQuote, false, true)
		assert.Equal(t, tt.graphic, out)
	}
}

type quoteRuneTest struct {
	in      rune
	out     string
	ascii   string
	graphic string
}

var quoterunetests = []quoteRuneTest{
	{'a', `a`, `a`, `a`},
	{'\a', `\a`, `\a`, `\a`},
	{'\\', `\\`, `\\`, `\\`},
	{0xFF, `ÿ`, `\u00ff`, `ÿ`},
	{0x263a, `☺`, `\u263a`, `☺`},
	{0xdead, `�`, `\ufffd`, `�`},
	{0xfffd, `�`, `\ufffd`, `�`},
	{0x0010ffff, `\U0010ffff`, `\U0010ffff`, `\U0010ffff`},
	{0x0010ffff + 1, `�`, `\ufffd`, `�`},
	{0x04, `\x04`, `\x04`, `\x04`},
	// Some differences between graphic and printable. Note the last column is double-quoted.
	{'\u00a0', `\u00a0`, `\u00a0`, "\u00a0"},
	{'\u2000', `\u2000`, `\u2000`, "\u2000"},
	{'\u3000', `\u3000`, `\u3000`, "\u3000"},
}

func TestEscapeFromGoRune(t *testing.T) {
	t.Parallel()

	for _, tt := range quoterunetests {
		out := token.EscapeRune(t.Context(), nil, tt.in, token.SingleQuote, false, false)
		assert.Equal(t, tt.out, string(out))
	}
}

func TestEscapeFromGoRuneASCII(t *testing.T) {
	t.Parallel()

	for _, tt := range quoterunetests {
		out := token.EscapeRune(t.Context(), nil, tt.in, token.SingleQuote, true, false)
		assert.Equal(t, tt.ascii, string(out))
	}
}

func TestEscapeFromGoRuneGraphic(t *testing.T) {
	t.Parallel()

	for _, tt := range quoterunetests {
		out := token.EscapeRune(t.Context(), nil, tt.in, token.SingleQuote, false, true)
		assert.Equal(t, tt.graphic, string(out))
	}
}
