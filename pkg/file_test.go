package pkg_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	pkg "github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"
)

func TestParseAndRoundTripSaveLoad(t *testing.T) {
	t.Parallel()

	doc, err := pkg.Parse(context.Background(), strings.NewReader("A=1\n"), "input.env")
	if err != nil {
		t.Fatalf("expected Parse to succeed, got %v", err)
	}

	filename := filepath.Join(t.TempDir(), "output.env")
	if err = pkg.Save(context.Background(), filename, doc); err != nil {
		t.Fatalf("expected Save to succeed, got %v", err)
	}

	loaded, err := pkg.Load(context.Background(), filename)
	if err != nil {
		t.Fatalf("expected Load to succeed, got %v", err)
	}

	if !loaded.Has("A") {
		t.Fatal("expected loaded document to contain A")
	}
}

func TestLoadMissingFileReturnsError(t *testing.T) {
	t.Parallel()

	_, err := pkg.Load(context.Background(), filepath.Join(t.TempDir(), "missing.env"))
	if err == nil {
		t.Fatal("expected error when loading missing file")
	}
}

func TestSaveEmptyDocumentReturnsError(t *testing.T) {
	t.Parallel()

	err := pkg.Save(context.Background(), filepath.Join(t.TempDir(), "empty.env"), ast.NewDocument())
	if err == nil {
		t.Fatal("expected Save to fail for empty document")
	}
}
