package parser_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/parser"
	"github.com/jippi/dottie/pkg/scanner"
	"github.com/jippi/dottie/pkg/token"

	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	t.Run("parse assigment successful", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected ast.Statement
		}{
			{
				name:  "unquoted value",
				input: "name=value",
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "name",
							Literal:      "value",
							Interpolated: "value",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Complete: true,
							Active:   true,
							Quote:    token.NoQuotes,
						},
					},
				},
			},
			{
				name:  "double quoted value",
				input: `name="value"`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "name",
							Literal:      "value",
							Interpolated: "value",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Complete: true,
							Active:   true,
							Quote:    token.DoubleQuotes,
						},
					},
				},
			},
			{
				name:  "single quoted value",
				input: `name='value'`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "name",
							Literal:      "value",
							Interpolated: "value",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote:    token.SingleQuotes,
							Complete: true,
							Active:   true,
						},
					},
				},
			},
			{
				name:  "name with assign and empty value",
				input: "name=",
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "name",
							Literal:      "",
							Interpolated: "",
							Active:       true,
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote: token.NoQuotes,
						},
					},
				},
			},
			{
				name:  "name without value",
				input: "name",
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "name",
							Literal:      "",
							Interpolated: "",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote:  token.NoQuotes,
							Active: true,
						},
					},
				},
			},
			{
				name:  "variable with blank lines",
				input: "\n\n\n\nname=\n\n\n",
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Newline{
							Blank: true,
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
						},
						&ast.Assignment{
							Name:         "name",
							Literal:      "",
							Interpolated: "",
							Position: ast.Position{
								Line:      5,
								FirstLine: 5,
								LastLine:  5,
							},
							Complete: false,
							Active:   true,
							Quote:    token.NoQuotes,
						},
						&ast.Newline{
							Blank: true,
							Position: ast.Position{
								Line:      6,
								FirstLine: 6,
								LastLine:  6,
							},
						},
					},
				},
			},
			{
				name:  "multiple variables",
				input: "DEBUG_HTTP_ADDR=:9090\nDEBUG_HTTP_IDLE_TIMEOUT=0s\nJAEGER_AGENT_ENDPOINT=jaeger-otlp-agent:6831",
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "DEBUG_HTTP_ADDR",
							Literal:      ":9090",
							Interpolated: ":9090",
							Complete:     true,
							Active:       true,
							Quote:        token.NoQuotes,
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
						},
						&ast.Assignment{
							Name:         "DEBUG_HTTP_IDLE_TIMEOUT",
							Literal:      "0s",
							Interpolated: "0s",
							Complete:     true,
							Active:       true,
							Quote:        token.NoQuotes,
							Position: ast.Position{
								Line:      2,
								FirstLine: 2,
								LastLine:  2,
							},
						},
						&ast.Assignment{
							Name:         "JAEGER_AGENT_ENDPOINT",
							Literal:      "jaeger-otlp-agent:6831",
							Interpolated: "jaeger-otlp-agent:6831",
							Complete:     true,
							Active:       true,
							Quote:        token.NoQuotes,
							Position: ast.Position{
								Line:      3,
								FirstLine: 3,
								LastLine:  3,
							},
						},
					},
				},
			},
			{
				name:  "variable with comments",
				input: "# comment 1\nDEBUG_HTTP_ADDR=:9090\n# comment 2",
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "DEBUG_HTTP_ADDR",
							Literal:      ":9090",
							Interpolated: ":9090",
							Position: ast.Position{
								Line:      2,
								FirstLine: 1,
								LastLine:  2,
							},
							Quote:    token.NoQuotes,
							Complete: true,
							Active:   true,
							Comments: []*ast.Comment{
								{
									Value: "# comment 1",
									Position: ast.Position{
										Line:      1,
										FirstLine: 1,
										LastLine:  1,
									},
								},
							},
						},
						&ast.Comment{
							Value: "# comment 2",
							Position: ast.Position{
								Line:      3,
								FirstLine: 3,
								LastLine:  3,
							},
						},
					},
				},
			},
			{
				name:  "newlines in quoted strings",
				input: `FOO="bar\nbaz"`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "FOO",
							Literal:      "bar\nbaz",
							Interpolated: "bar\nbaz",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Complete: true,
							Active:   true,
							Quote:    token.DoubleQuotes,
						},
					},
				},
			},
			{
				name:  "newlines in naked strings",
				input: `FOO=bar\nbaz`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "FOO",
							Literal:      "bar\nbaz",
							Interpolated: "bar\nbaz",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote:    token.NoQuotes,
							Complete: true,
							Active:   true,
						},
					},
				},
			},
			{
				name:  "single quotes inside double quotes",
				input: `FOO="'d'"`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "FOO",
							Literal:      "'d'",
							Interpolated: "'d'",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote:    token.DoubleQuotes,
							Complete: true,
							Active:   true,
						},
					},
				},
			},
			{
				name:  `variable with several "=" in the value`,
				input: `FOO=foobar=`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "FOO",
							Literal:      "foobar=",
							Interpolated: "foobar=",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote:    token.NoQuotes,
							Complete: true,
							Active:   true,
						},
					},
				},
			},
			{
				name:  `inline comments is a part of value`,
				input: `FOO=bar # this is foo`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "FOO",
							Literal:      "bar # this is foo",
							Interpolated: "bar # this is foo",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Complete: true,
							Active:   true,
							Quote:    token.NoQuotes,
						},
					},
				},
			},
			{
				name:  `allows # in double quoted value`,
				input: `FOO="bar#baz"`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "FOO",
							Literal:      "bar#baz",
							Interpolated: "bar#baz",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote:    token.DoubleQuotes,
							Complete: true,
							Active:   true,
						},
					},
				},
			},
			{
				name:  `allows # in single quoted value`,
				input: `FOO='bar#baz'`,
				expected: &ast.Document{
					Statements: []ast.Statement{
						&ast.Assignment{
							Name:         "FOO",
							Literal:      "bar#baz",
							Interpolated: "bar#baz",
							Position: ast.Position{
								Line:      1,
								FirstLine: 1,
								LastLine:  1,
							},
							Quote:    token.SingleQuotes,
							Complete: true,
							Active:   true,
						},
					},
				},
			},
		}

		for _, tt := range tests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				s := scanner.New(tt.input)
				p := parser.New(s, "-")

				stmts, err := p.Parse()
				require.NoError(t, err)
				require.EqualValues(t, tt.expected, stmts)
			})
		}
	})

	t.Run("returns error on invalid input", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "escaped double quotes",
				input: `FOO="escaped\"bar"`,
			},
			{
				name:  "value with space after equal sign",
				input: `FOO= bar`,
			},
			{
				name:  "value with space before equal sign",
				input: `FOO =bar`,
			},
			{
				name:  "leading tab",
				input: "\tFOO=bar",
			},
			{
				name:  "leading whitespace",
				input: "  FOO=bar",
			},
		}

		for _, tt := range tests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				s := scanner.New(tt.input)
				p := parser.New(s, "-")

				stmts, err := p.Parse()
				require.Error(t, err, "expected an error")
				require.Nil(t, stmts, "did not expect a statement")
			})
		}
	})
}
