package upsert_test

import (
	"strings"
	"testing"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/ast/upsert"
)

func TestPlacementAndSettingPublicBehavior(t *testing.T) {
	t.Parallel()

	if !upsert.AddAfterKey.RequiresKey() || !upsert.AddBeforeKey.RequiresKey() {
		t.Fatal("expected relative placements to require key")
	}

	if upsert.AddLast.RequiresKey() || upsert.AddFirst.RequiresKey() {
		t.Fatal("expected non-relative placements not to require key")
	}

	combined := upsert.SkipIfSame | upsert.Validate
	if !combined.Has(upsert.SkipIfSame) || !combined.Has(upsert.Validate) {
		t.Fatal("expected combined settings to include SkipIfSame and Validate")
	}

	if got := upsert.AddLast.String(); got != "Placement<AddLast>" {
		t.Fatalf("unexpected AddLast string: %q", got)
	}
}

func TestSkippedStatementErrorMessage(t *testing.T) {
	t.Parallel()

	err := upsert.SkippedStatementError{Key: "APP_PORT", Reason: "already exists"}
	if got := err.Error(); got != "Key [ APP_PORT ] was skipped: already exists" {
		t.Fatalf("unexpected skipped error text: %q", got)
	}
}

func TestPlacementOptionValidation(t *testing.T) {
	t.Parallel()

	document := ast.NewDocument()

	_, err := upsert.New(document, upsert.WithPlacement(upsert.AddAfterKey))
	if err == nil || !strings.Contains(err.Error(), "does requires a KEY") {
		t.Fatalf("expected placement error for key-required mode, got %v", err)
	}

	_, err = upsert.New(document, upsert.WithPlacementRelativeToKey(upsert.AddBeforeKey, "MISSING_KEY"))
	if err == nil || !strings.Contains(err.Error(), "does not exists") {
		t.Fatalf("expected missing-key placement error, got %v", err)
	}
}

func TestPlacementIgnoringEmptyIsNoop(t *testing.T) {
	t.Parallel()

	document := ast.NewDocument()

	created, err := upsert.New(document, upsert.WithPlacementIgnoringEmpty(upsert.AddBeforeKey, ""))
	if err != nil {
		t.Fatalf("expected empty key placement option to be no-op, got %v", err)
	}

	if created == nil {
		t.Fatal("expected upserter to be created")
	}
}
