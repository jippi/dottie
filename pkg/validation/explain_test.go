package validation_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

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
