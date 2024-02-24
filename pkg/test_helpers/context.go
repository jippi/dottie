package test_helpers

import (
	"bytes"
	"context"
	"testing"

	"github.com/jippi/dottie/pkg/tui"
	"github.com/neilotoole/slogt"
	slogctx "github.com/veqryn/slog-context"
)

func CreateContext(t *testing.T, out, err *bytes.Buffer) context.Context {
	t.Helper()

	if out == nil {
		out = &bytes.Buffer{}
	}

	if err == nil {
		err = &bytes.Buffer{}
	}

	ctx := tui.NewContextWithoutLogger(context.Background(), out, err)
	// return ctx

	return slogctx.NewCtx(ctx, slogt.New(t, slogt.Text()))
}
