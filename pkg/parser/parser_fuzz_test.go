package parser_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/parser"
	"github.com/jippi/dottie/pkg/scanner"
)

func FuzzParseDoesNotPanic(f *testing.F) {
	seeds := []string{
		"A=B",
		"# comment\nA=$B",
		"A=\"line1\\nline2\"",
		"# @dottie/validate required\nA=",
		"################################################################################\n# group\n################################################################################\nA=1",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ctx := t.Context()

		scan := scanner.New(input)

		doc, err := parser.New(ctx, scan, "-").Parse(ctx)
		if err != nil {
			return
		}

		if doc == nil {
			t.Fatal("expected document when parse succeeds")
		}

		_ = doc.InterpolateAll(ctx)
	})
}
