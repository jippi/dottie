package ast_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/token"
)

func TestInterpolateAllDetectsCyclicDependencies(t *testing.T) {
	t.Parallel()

	doc := ast.NewDocument()
	doc.Statements = []ast.Statement{
		&ast.Assignment{
			Name:         "A",
			Literal:      "$B",
			Interpolated: "$B",
			Enabled:      true,
			Quote:        token.NoQuote,
			Position: ast.Position{
				File:  "test.env",
				Line:  1,
				Index: 0,
			},
		},
		&ast.Assignment{
			Name:         "B",
			Literal:      "$A",
			Interpolated: "$A",
			Enabled:      true,
			Quote:        token.NoQuote,
			Position: ast.Position{
				File:  "test.env",
				Line:  2,
				Index: 1,
			},
		},
	}

	doc.Initialize(context.Background())

	err := doc.InterpolateAll(context.Background())
	if err == nil {
		t.Fatal("expected cycle interpolation error, got nil")
	}

	if !strings.Contains(err.Error(), "cyclic dependency detected") {
		t.Fatalf("expected cycle error message, got %q", err.Error())
	}
}

func TestInitializeRebuildsDependentsFromCurrentDependencies(t *testing.T) {
	t.Parallel()

	base := &ast.Assignment{
		Name:         "BASE",
		Literal:      "value",
		Interpolated: "value",
		Enabled:      true,
		Quote:        token.NoQuote,
		Position: ast.Position{
			File:  "test.env",
			Line:  1,
			Index: 0,
		},
	}

	consumer := &ast.Assignment{
		Name:         "CONSUMER",
		Literal:      "$BASE",
		Interpolated: "$BASE",
		Enabled:      true,
		Quote:        token.NoQuote,
		Position: ast.Position{
			File:  "test.env",
			Line:  2,
			Index: 1,
		},
	}

	doc := ast.NewDocument()
	doc.Statements = []ast.Statement{base, consumer}

	doc.Initialize(context.Background())

	if base.Dependents == nil || base.Dependents["CONSUMER"] == nil {
		t.Fatal("expected BASE to have CONSUMER as dependent after initial initialize")
	}

	consumer.Literal = "plain"

	doc.Initialize(context.Background())

	if len(base.Dependents) != 0 {
		t.Fatalf("expected BASE dependents to be rebuilt and emptied, got %+v", base.Dependents)
	}
}

func TestInterpolationMapperDoesNotResolveDisabledAssignments(t *testing.T) {
	t.Parallel()

	disabled := &ast.Assignment{
		Name:         "SECRET",
		Literal:      "secret-value",
		Interpolated: "secret-value",
		Enabled:      false,
		Quote:        token.NoQuote,
		Position: ast.Position{
			File:  "test.env",
			Line:  1,
			Index: 0,
		},
	}

	target := &ast.Assignment{
		Name:         "PUBLIC",
		Literal:      "$SECRET",
		Interpolated: "$SECRET",
		Enabled:      true,
		Quote:        token.NoQuote,
		Position: ast.Position{
			File:  "test.env",
			Line:  2,
			Index: 1,
		},
	}

	doc := ast.NewDocument()
	doc.Statements = []ast.Statement{disabled, target}

	doc.Initialize(context.Background())

	_, ok := doc.InterpolationMapper(target)("SECRET")
	if ok {
		t.Fatal("expected disabled assignment to be inaccessible via interpolation mapper")
	}
}
