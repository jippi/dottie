package scanner_test

import (
	"strings"
	"testing"

	"github.com/jippi/dottie/pkg/scanner"
	"github.com/jippi/dottie/pkg/token"

	"github.com/stretchr/testify/assert"
)

const (
	oversizedValuePadding = 70 * 1024
	oversizedInputPadding = 200 * 1024
)

func TestScanner_NextToken_Trivial(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		input             string
		expectedTokenType token.Type
		expectedLiteral   string
	}{
		{
			name:              "empty input",
			input:             "",
			expectedTokenType: token.EOF,
			expectedLiteral:   token.EOF.String(),
		},
		{
			name:              "BOM",
			input:             "\uFEFF",
			expectedTokenType: token.EOF,
			expectedLiteral:   token.EOF.String(),
		},
		{
			name:              "BOM and then assignment",
			input:             "\uFEFF=",
			expectedTokenType: token.Assign,
			expectedLiteral:   token.Assign.String(),
		},
		{
			name:              "new lines",
			input:             "\n\n",
			expectedTokenType: token.NewLine,
			expectedLiteral:   "\n",
		},
		{
			name:              "assignment",
			input:             "=",
			expectedTokenType: token.Assign,
			expectedLiteral:   token.Assign.String(),
		},
		{
			name:              "space",
			input:             " ",
			expectedTokenType: token.Space,
			expectedLiteral:   " ",
		},
		{
			name:              "tab",
			input:             "\t",
			expectedTokenType: token.Space,
			expectedLiteral:   "\t", // `	`
		},
		{
			name:              "vertical tab",
			input:             "\v",
			expectedTokenType: token.Space,
			expectedLiteral:   "\v", // ``
		},
		{
			name:              "form feed",
			input:             "\f",
			expectedTokenType: token.Space,
			expectedLiteral:   "\f", // ``
		},
		{
			name:              "carriage return",
			input:             "\r",
			expectedTokenType: token.Space,
			expectedLiteral:   "\r", // `␍`
		},
		{
			name:              "comment",
			input:             "# comment",
			expectedTokenType: token.Comment,
			expectedLiteral:   "# comment",
		},
		{
			name:              "comment with a new line",
			input:             "# comment\n",
			expectedTokenType: token.Comment,
			expectedLiteral:   "# comment",
		},
		{
			name:              "identifier",
			input:             "valid.identifier='valid value'",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "valid.identifier",
		},
		{
			name:              "double quoted value",
			input:             `"valid value"`,
			expectedTokenType: token.Value,
			expectedLiteral:   "valid value",
		},
		{
			name:              "double quoted value with escaped new line",
			input:             `"valid value\n"`,
			expectedTokenType: token.Value,
			expectedLiteral:   "valid value\n",
		},
		{
			name:              "single quoted value",
			input:             `'valid value'`,
			expectedTokenType: token.RawValue,
			expectedLiteral:   `valid value`,
		},
		{
			name:              "single quoted value with escaped new line",
			input:             `'valid value \n'`,
			expectedTokenType: token.RawValue,
			expectedLiteral:   `valid value \n`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := scanner.New(tt.input)

			actual := sc.NextToken(t.Context())
			assert.Equal(t, tt.expectedTokenType, actual.Type)
			assert.Equal(t, tt.expectedLiteral, actual.Literal)
		})
	}
}

func TestScanner_NextToken_Valid_Identifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		input             string
		expectedTokenType token.Type
		expectedLiteral   string
	}{
		{
			name:              "ASCII lower letters",
			input:             "abcdefghijklmnopqrstuvwxyz",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "abcdefghijklmnopqrstuvwxyz",
		},
		{
			name:              "ASCII upper letters",
			input:             "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			name:              "ASCII mixed letters",
			input:             "ABCDEFGHIJKLMnopqrstuvwxyz",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "ABCDEFGHIJKLMnopqrstuvwxyz",
		},
		{
			name:              "digits",
			input:             "1234567890",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "1234567890",
		},
		{
			name:              "digits + ASCII letters",
			input:             "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			name:              "underscore",
			input:             "_",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "_",
		},
		{
			name:              "multiple underscores",
			input:             "___",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "___",
		},
		{
			name:              "all special symbols",
			input:             ".-,_",
			expectedTokenType: token.Identifier,
			expectedLiteral:   ".-,_",
		},
		{
			name:              "cyrillic lower letters",
			input:             "абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
		},
		{
			name:              "cyrillic upper letters",
			input:             "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ",
		},
		{
			name:              "chinese letters",
			input:             "一个类型",
			expectedTokenType: token.Identifier,
			expectedLiteral:   "一个类型",
		},
		{
			name:              "mixed quotes should fail",
			input:             `'valid value \n"`,
			expectedTokenType: token.Illegal,
			expectedLiteral:   `valid value \n"`,
		},
		{
			name:              "multiline single quote with double quote inside OK",
			input:             `'valid value \n"'`,
			expectedTokenType: token.RawValue,
			expectedLiteral:   `valid value \n"`,
		},
		{
			name:              "multiline double quote with single quote inside OK",
			input:             `"valid value \n'"`,
			expectedTokenType: token.Value,
			expectedLiteral:   "valid value \n'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := scanner.New(tt.input)

			actual := sc.NextToken(t.Context())
			assert.Equal(t, tt.expectedTokenType.String(), actual.Type.String())
			assert.Equal(t, tt.expectedLiteral, actual.Literal)
		})
	}
}

func TestScanner_NextToken_OversizedUnquotedValueReturnsIllegal(t *testing.T) {
	t.Parallel()

	input := "A=" + strings.Repeat("\x8f", oversizedValuePadding)
	sc := scanner.New(input)

	assert.Equal(t, token.Identifier, sc.NextToken(t.Context()).Type)
	assert.Equal(t, token.Assign, sc.NextToken(t.Context()).Type)

	value := sc.NextToken(t.Context())
	assert.Equal(t, token.Illegal, value.Type)
	assert.Equal(t, "value exceeds maximum supported length", value.Literal)
}

func TestScanner_NextToken_OversizedQuotedValueReturnsIllegal(t *testing.T) {
	t.Parallel()

	input := "A=\"" + strings.Repeat("a", oversizedValuePadding) + "\""
	sc := scanner.New(input)

	assert.Equal(t, token.Identifier, sc.NextToken(t.Context()).Type)
	assert.Equal(t, token.Assign, sc.NextToken(t.Context()).Type)

	value := sc.NextToken(t.Context())
	assert.Equal(t, token.Illegal, value.Type)
}

func TestScanner_NextToken_OversizedInputIsCapped(t *testing.T) {
	t.Parallel()

	input := "# " + strings.Repeat("a", oversizedInputPadding)
	sc := scanner.New(input)

	comment := sc.NextToken(t.Context())
	assert.Equal(t, token.Comment, comment.Type)
	assert.Less(t, len(comment.Literal), len(input))
}

func TestScanner_NextToken_Naked_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		input             string
		expectedTokenType token.Type
		expectedLiteral   string
	}{
		{
			name:              "empty input after =",
			input:             "=",
			expectedTokenType: token.EOF,
			expectedLiteral:   token.EOF.String(),
		},
		{
			name:              "new lines after =",
			input:             "=\n",
			expectedTokenType: token.NewLine,
			expectedLiteral:   "\n",
		},
		{
			name:              "valid assignment",
			input:             "=valid value",
			expectedTokenType: token.Value,
			expectedLiteral:   `valid value`,
		},
		{
			name:              "valid assignment with escaped new line",
			input:             `=valid value \n`,
			expectedTokenType: token.Value,
			expectedLiteral:   "valid value \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := scanner.New(tt.input)
			assign := sc.NextToken(t.Context())
			assert.Equal(t, token.Assign, assign.Type)
			assert.Equal(t, token.Assign.String(), assign.Literal)

			actual := sc.NextToken(t.Context())
			assert.Equal(t, tt.expectedTokenType.String(), actual.Type.String())
			assert.Equal(t, tt.expectedLiteral, actual.Literal)
		})
	}
}

func TestScanner_NextToken_InvalidUTF8DoesNotPanic(t *testing.T) {
	t.Parallel()

	sc := scanner.New(string([]byte{0xff}))

	actual := sc.NextToken(t.Context())
	assert.Equal(t, token.Illegal, actual.Type)
}

func TestScanner_NextToken_Illegal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{
			name:            "illegal identifier",
			input:           "$invalid.name$=value",
			expectedLiteral: "$",
		},
		{
			name:            "not-paired double quotes",
			input:           `"quotes must be closed`,
			expectedLiteral: "quotes must be closed",
		},
		{
			name:            "not-paired mixed quotes",
			input:           `"quotes should not be mixed'`,
			expectedLiteral: "quotes should not be mixed'",
		},
		{
			name: "not-paired double quotes with new line",
			input: `"quotes must be closed
`,
			expectedLiteral: "quotes must be closed\n",
		},
		{
			name:            "not-paired single quotes",
			input:           `'quotes must be closed`,
			expectedLiteral: "quotes must be closed",
		},
		{
			name: "not-paired single quotes with new line",
			input: `'quotes must be closed
`,
			expectedLiteral: "quotes must be closed\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := scanner.New(tt.input)

			actual := sc.NextToken(t.Context())
			assert.Equal(t, token.Illegal.String(), actual.Type.String())
			assert.Equal(t, tt.expectedLiteral, actual.Literal)
		})
	}
}

func TestScanner_NextToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name: "illegal value quote double quotes",
			input: `x="yxc
`,
			expected: []token.Token{
				{Type: token.Identifier, Literal: "x"},
				{Type: token.Assign, Literal: token.Assign.String()},
				{Type: token.Illegal, Literal: "yxc\n"},
				{Type: token.EOF, Literal: token.EOF.String()},
			},
		},
		{
			name: "illegal value quote single quotes",
			input: `x='yxc
`,
			expected: []token.Token{
				{Type: token.Identifier, Literal: "x"},
				{Type: token.Assign, Literal: token.Assign.String()},
				{Type: token.Illegal, Literal: "yxc\n"},
				{Type: token.EOF, Literal: token.EOF.String()},
			},
		},
		{
			name: "naked value",
			input: `x=yxc
`,
			expected: []token.Token{
				{Type: token.Identifier, Literal: "x"},
				{Type: token.Assign, Literal: token.Assign.String()},
				{Type: token.Value, Literal: "yxc"},
				{Type: token.NewLine, Literal: "\n"},
				{Type: token.EOF, Literal: token.EOF.String()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := scanner.New(tt.input)

			counter := 0

			for {
				actual := sc.NextToken(t.Context())
				expected := tt.expected[counter]

				assert.Equal(t, expected.Type, actual.Type)
				assert.Equal(t, expected.Literal, actual.Literal)

				if actual.Type == token.EOF {
					break
				}

				counter++
			}
		})
	}
}
