package validation_test

import (
	"bytes"
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
	"github.com/jippi/dottie/pkg/validation"
)

func TestExplainWithErrorIncludesMessage(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := tui.NewContextWithoutLogger(context.Background(), &stdout, &stderr)
	result := validation.Explain(ctx, nil, errors.New("boom"), nil, false, false)

	if !strings.Contains(result, "boom") {
		t.Fatalf("expected Explain output to contain error message, got %q", result)
	}
}

func TestExplainColorizesHighlightedValidationParts(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "1")
	t.Setenv("TERM", "xterm-256color")

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := tui.NewContextWithoutLogger(context.Background(), &stdout, &stderr)

	err := validator.New().Var("bad", "oneof=a b")
	if err == nil {
		t.Fatal("expected validation error")
	}

	assignment := &ast.Assignment{
		Name:         "TEST_KEY",
		Interpolated: "bad",
	}

	result := validation.Explain(ctx, nil, err, assignment, false, false)

	ansiPattern := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	if !ansiPattern.MatchString(result) {
		t.Fatalf("expected colorized output with ANSI escapes, got %q", result)
	}

	if !regexp.MustCompile(`\[\x1b\[[0-9;]*mbad\x1b\[[0-9;]*m\]`).MatchString(result) {
		t.Fatalf("expected highlighted value segment in output, got %q", result)
	}

	if !regexp.MustCompile(`\[\x1b\[[0-9;]*ma b\x1b\[[0-9;]*m\]`).MatchString(result) {
		t.Fatalf("expected highlighted rule parameter segment in output, got %q", result)
	}
}
