package token_test

import (
	"context"
	"fmt"
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual, err := token.Unescape(context.TODO(), tt.input, token.DoubleQuote)
			require.NoError(t, err)

			require.Equal(t, tt.expected, actual)
		})
	}
}

var unescapeTests = []struct {
	in    string
	out   string
	quote token.Quote
}{
	{``, "", token.DoubleQuote},                                     // 0
	{`a`, "a", token.DoubleQuote},                                   // 1
	{`abc`, "abc", token.DoubleQuote},                               // 2
	{`☺`, "☺", token.DoubleQuote},                                   // 3
	{`hello world`, "hello world", token.DoubleQuote},               // 4
	{`\xFF`, "\xFF", token.DoubleQuote},                             // 5
	{`\377`, "\377", token.DoubleQuote},                             // 6
	{`\u1234`, "\u1234", token.DoubleQuote},                         // 7
	{`\U00010111`, "\U00010111", token.DoubleQuote},                 // 8
	{`\U0001011111`, "\U0001011111", token.DoubleQuote},             // 9
	{`\a\b\f\n\r\t\v\\\"`, "\a\b\f\n\r\t\v\\\"", token.DoubleQuote}, // 10
	{`'`, "'", token.DoubleQuote},                                   // 11
	{`a`, "a", token.DoubleQuote},                                   // 12
	{`☹`, "☹", token.DoubleQuote},                                   // 13
	{`\a`, "\a", token.DoubleQuote},                                 // 14
	{`\x10`, "\x10", token.DoubleQuote},                             // 15
	{`\377`, "\377", token.DoubleQuote},                             // 16
	{`\u1234`, "\u1234", token.DoubleQuote},                         // 17
	{`\U00010111`, "\U00010111", token.DoubleQuote},                 // 18
	{`\t`, "\t", token.DoubleQuote},                                 // 19
	{` `, " ", token.DoubleQuote},                                   // 20
	{`\'`, "'", token.DoubleQuote},                                  // 21
	{`"`, "\"", token.DoubleQuote},                                  // 22
	{"", ``, token.DoubleQuote},                                     // 23
	{"a", `a`, token.DoubleQuote},                                   // 24
	{"abc", `abc`, token.DoubleQuote},                               // 25
	{"☺", `☺`, token.DoubleQuote},                                   // 26
	{"hello world", `hello world`, token.DoubleQuote},               // 27
	{"\\", `\`, token.DoubleQuote},                                  // 28
	{"\n", "\n", token.DoubleQuote},                                 // 29
	{"	", `	`, token.DoubleQuote},                                   // 30
	{" ", ` `, token.DoubleQuote},                                   // 31
}

func TestUnquote(t *testing.T) {
	t.Parallel()

	for idx, tt := range unescapeTests {
		t.Run(fmt.Sprintf("unquote-tests-%d", idx), func(t *testing.T) {
			t.Parallel()

			testUnescapeHelper(t, tt.in, tt.out, tt.quote)
		})
	}

	for idx, tt := range quotetests {
		t.Run(fmt.Sprintf("quote-tests-%d", idx), func(t *testing.T) {
			t.Parallel()

			testUnescapeHelper(t, tt.out, tt.in, token.DoubleQuote)
		})
	}
}

// Issue 23685: invalid UTF-8 should not go through the fast path.
func TestUnescapeInvalidUTF8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: `foo`, want: "foo"},
		{in: `` + "\xc0" + ``, want: "\xef\xbf\xbd"},
		{in: `a` + "\xc0" + ``, want: "a\xef\xbf\xbd"},
		{in: `\t` + "\xc0" + ``, want: "\t\xef\xbf\xbd"},
	}

	for _, tt := range tests {
		testUnescapeHelper(t, tt.in, tt.want, token.DoubleQuote)
	}
}

func testUnescapeHelper(t *testing.T, in, want string, quote token.Quote) {
	t.Helper()

	got, err := token.Unescape(context.TODO(), in, quote)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
