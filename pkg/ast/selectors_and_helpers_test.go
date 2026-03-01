package ast_test

import (
	"errors"
	"testing"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/token"
)

func TestExcludeComments(t *testing.T) {
	t.Parallel()

	comment := &ast.Comment{Value: "# hello"}
	assignment := &ast.Assignment{Name: "A"}

	if got := ast.ExcludeComments(comment); got != ast.Exclude {
		t.Fatalf("expected comment to be excluded, got %v", got)
	}

	if got := ast.ExcludeComments(assignment); got != ast.Keep {
		t.Fatalf("expected assignment to be kept, got %v", got)
	}
}

func TestAssignmentStateSelectors(t *testing.T) {
	t.Parallel()

	enabled := &ast.Assignment{Enabled: true}
	disabled := &ast.Assignment{Enabled: false}

	if got := ast.ExcludeActiveAssignments(enabled); got != ast.Exclude {
		t.Fatalf("expected enabled assignment to be excluded by ExcludeActiveAssignments, got %v", got)
	}

	if got := ast.ExcludeDisabledAssignments(disabled); got != ast.Exclude {
		t.Fatalf("expected disabled assignment to be excluded by ExcludeDisabledAssignments, got %v", got)
	}

	if got := ast.ExcludeActiveAssignments(disabled); got != ast.Keep {
		t.Fatalf("expected disabled assignment to be kept by ExcludeActiveAssignments, got %v", got)
	}

	if got := ast.ExcludeDisabledAssignments(enabled); got != ast.Keep {
		t.Fatalf("expected enabled assignment to be kept by ExcludeDisabledAssignments, got %v", got)
	}
}

func TestExcludeHiddenViaAnnotation(t *testing.T) {
	t.Parallel()

	hidden := &ast.Assignment{
		Comments: []*ast.Comment{{
			Annotation: &token.Annotation{Key: "dottie/hidden"},
		}},
	}

	notHidden := &ast.Assignment{}

	if got := ast.ExcludeHiddenViaAnnotation(hidden); got != ast.Exclude {
		t.Fatalf("expected hidden assignment to be excluded, got %v", got)
	}

	if got := ast.ExcludeHiddenViaAnnotation(notHidden); got != ast.Keep {
		t.Fatalf("expected non-hidden assignment to be kept, got %v", got)
	}
}

func TestKeySelectors(t *testing.T) {
	t.Parallel()

	assignment := &ast.Assignment{Name: "APP_PORT"}

	if got := ast.RetainKeyPrefix("APP_")(assignment); got != ast.Keep {
		t.Fatalf("expected RetainKeyPrefix to keep APP_PORT, got %v", got)
	}

	if got := ast.RetainKeyPrefix("DB_")(assignment); got != ast.Exclude {
		t.Fatalf("expected RetainKeyPrefix to exclude APP_PORT for DB_ prefix, got %v", got)
	}

	if got := ast.RetainExactKey("APP_PORT", "APP_HOST")(assignment); got != ast.Keep {
		t.Fatalf("expected RetainExactKey to keep APP_PORT, got %v", got)
	}

	if got := ast.RetainExactKey("DB_PORT")(assignment); got != ast.Exclude {
		t.Fatalf("expected RetainExactKey to exclude APP_PORT, got %v", got)
	}

	if got := ast.ExcludeKeyPrefix("APP_")(assignment); got != ast.Exclude {
		t.Fatalf("expected ExcludeKeyPrefix to exclude APP_PORT for APP_ prefix, got %v", got)
	}

	if got := ast.ExcludeKeyPrefix("DB_")(assignment); got != ast.Keep {
		t.Fatalf("expected ExcludeKeyPrefix to keep APP_PORT for DB_ prefix, got %v", got)
	}
}

func TestRetainGroupSelector(t *testing.T) {
	t.Parallel()

	group := &ast.Group{Name: "# app"}
	assignment := &ast.Assignment{Name: "APP_PORT", Group: group}
	comment := &ast.Comment{Value: "# docs", Group: group}

	selector := ast.RetainGroup("app")

	if got := selector(group); got != ast.Keep {
		t.Fatalf("expected group to be kept, got %v", got)
	}

	if got := selector(assignment); got != ast.Keep {
		t.Fatalf("expected assignment in matching group to be kept, got %v", got)
	}

	if got := selector(comment); got != ast.Keep {
		t.Fatalf("expected comment in matching group to be kept, got %v", got)
	}

	other := &ast.Assignment{Name: "DB_PORT", Group: &ast.Group{Name: "# db"}}

	if got := selector(other); got != ast.Exclude {
		t.Fatalf("expected assignment in non-matching group to be excluded, got %v", got)
	}
}

func TestContextualError(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("boom")

	if got := ast.ContextualError(nil, nil); got != nil {
		t.Fatalf("expected nil when error is nil, got %v", got)
	}

	if got := ast.ContextualError(nil, baseErr); !errors.Is(got, baseErr) {
		t.Fatalf("expected original error when statement is nil, got %v", got)
	}

	assignment := &ast.Assignment{Position: ast.Position{File: "example.env", Line: 7}}

	got := ast.ContextualError(assignment, baseErr)
	if got == nil {
		t.Fatal("expected contextual error, got nil")
	}

	if !errors.Is(got, baseErr) {
		t.Fatalf("expected wrapped error to match original, got %v", got)
	}

	if got.Error() != "boom (example.env:7)" {
		t.Fatalf("unexpected contextual error message: %q", got.Error())
	}
}

func TestNewlineIsAndType(t *testing.T) {
	t.Parallel()

	newline := &ast.Newline{}

	if !newline.Is(&ast.Newline{}) {
		t.Fatal("expected Newline.Is to return true for same type")
	}

	if newline.Is(nil) {
		t.Fatal("expected Newline.Is to return false for nil")
	}

	if got := newline.Type(); got == "" {
		t.Fatal("expected Newline.Type to return non-empty type")
	}
}
