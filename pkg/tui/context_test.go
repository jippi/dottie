package tui_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/jippi/dottie/pkg/tui"
)

func TestNewContextWithoutLoggerStoresWriters(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := tui.NewContextWithoutLogger(context.Background(), &stdout, &stderr)
	outWriter := tui.StdoutFromContext(ctx).GetWriter()
	errWriter := tui.StderrFromContext(ctx).GetWriter()

	if outWriter != &stdout {
		t.Fatal("expected stdout writer to be stored in context")
	}

	if errWriter != &stderr {
		t.Fatal("expected stderr writer to be stored in context")
	}
}
